package util

import (
	"log"

	ort "github.com/yalue/onnxruntime_go"
)

// init 函数在程序启动时执行
func init() {
	log.Println("Model loading initialized.")
}

// check 函数：加载模型并进行推理
func Check(inputData []float32) ([]float32, error) {

	ort.SetSharedLibraryPath("E:/onnxruntime-1.20.2/onnxruntime-win-x64-1.20.0/onnxruntime-win-x64-1.20.0/lib/onnxruntime.dll")
	err := ort.InitializeEnvironment()
	if err != nil {
		panic(err)
	}
	defer ort.DestroyEnvironment()

	// 读取 ONNX 模型文件
	modelPath := "./model_64/modified_mobilenetv2_dst_64.onnx"

	// 创建输入张量，shape: (1, 1, 8, 8)
	inputShape := ort.NewShape(1, 1, 8, 8)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	defer inputTensor.Destroy()

	// 创建输出张量，shape: (1, 2)
	outputShape := ort.NewShape(1, 2)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	defer outputTensor.Destroy()

	// 创建 ONNX Runtime 运行会话
	session, err := ort.NewAdvancedSession(modelPath,
		[]string{"input"}, []string{"output"},
		[]ort.Value{inputTensor}, []ort.Value{outputTensor}, nil)
	defer session.Destroy()

	//运行
	err = session.Run()
	if err != nil {
		log.Fatalf("Model inference failed: %v", err)
	}

	result := outputTensor.GetData()

	return result, nil
}
