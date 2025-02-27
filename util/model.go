package util

import (
	"log"

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

// check 函数：加载模型并进行推理
func Check(inputData []float32) ([]float32, error) {
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
