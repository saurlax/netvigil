package util

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/spf13/viper"
	ort "github.com/yalue/onnxruntime_go"
)

func init() {
	libPath := viper.GetString("ortlib_path")

	if libPath != "" {
		log.Printf("Using ONNX Runtime shared library: %s", libPath)
		ort.SetSharedLibraryPath(libPath)
		err := ort.InitializeEnvironment()
		if err != nil {
			panic(err)
		}
	}
}

func Check(packet gopacket.Packet) bool {

	image := PreprocessPacket(packet)

	result, err := _check(image)
	if err != nil {
		log.Fatalf("Error during model inference: %v", err)
	}

	if result[0] > result[1] {
		return false
	} else {
		return true
	}

}

// packet转换为张量的函数
func PreprocessPacket(packet gopacket.Packet) []float32 {
	// 匿名化
	raw := AnonymizePacket(packet)

	// 保留前64字节，不足补0
	fixed := make([]byte, 64)
	copy(fixed, raw)

	// 转换为 float32 并归一化
	floatData := make([]float32, 64)
	for i := 0; i < 64; i++ {
		floatData[i] = float32(fixed[i]) / 255.0
	}

	// 标准化,由模型决定
	var mean float32
	var std float32
	mean = 0.2515
	std = 0.3778

	for i := range floatData {
		floatData[i] = (floatData[i] - mean) / std
	}

	return floatData
}

// 匿名化函数
func AnonymizePacket(packet gopacket.Packet) []byte {
	data := packet.Data()
	newData := make([]byte, len(data))
	copy(newData, data)

	// 匿名化 MAC 地址
	if len(newData) >= 14 {
		copy(newData[0:6], make([]byte, 6))  // dst MAC
		copy(newData[6:12], make([]byte, 6)) // src MAC
	}

	if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
		//处理ipv4
		ipHeaderOffset := 14                                                // Ethernet header
		copy(newData[ipHeaderOffset+12:ipHeaderOffset+16], make([]byte, 4)) // src IP
		copy(newData[ipHeaderOffset+16:ipHeaderOffset+20], make([]byte, 4)) // dst IP

	} else if ipv6Layer := packet.Layer(layers.LayerTypeIPv6); ipv6Layer != nil {
		//处理ipv6，  //好像输入的张量大小容不下ipv6
		ipHeaderOffset := 14                                                 // Ethernet header
		copy(newData[ipHeaderOffset+8:ipHeaderOffset+24], make([]byte, 16))  // src IP
		copy(newData[ipHeaderOffset+24:ipHeaderOffset+40], make([]byte, 16)) // dst IP
	}

	return newData
}

// check 函数：加载模型并进行推理
func _check(inputData []float32) ([]float32, error) {
	modelPath := viper.GetString("model_path")

	// 创建输入张量，shape: (1, 1, 8, 8)
	inputShape := ort.NewShape(1, 1, 8, 8)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	defer inputTensor.Destroy()

	// 创建输出张量，shape: (1, 2)
	outputShape := ort.NewShape(1, 2)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	defer outputTensor.Destroy()

	// 读取 ONNX 模型文件
	session, err := ort.NewAdvancedSession(modelPath,
		[]string{"input"}, []string{"output"}, []ort.Value{inputTensor}, []ort.Value{outputTensor}, nil)
	if err != nil {
		log.Fatalf("Error creating ONNX session: %v", err)
	}
	defer session.Destroy()

	//运行
	err = session.Run()
	if err != nil {
		log.Fatalf("Model inference failed: %v", err)
	}

	result := outputTensor.GetData()

	return result, nil
}
