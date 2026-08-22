package shader

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type ShaderProgram struct {
	FragmentShaderFilePath string
	VertexShaderFilePath   string
	programID              uint32
}

func (shp *ShaderProgram) Use() {
	gl.UseProgram(shp.programID)
}

func (shp *ShaderProgram) Delete() {
	gl.DeleteProgram(shp.programID)
}

func (shp *ShaderProgram) Compile() (error) {
	vertShader, err := compileShader(shp.VertexShaderFilePath, gl.VERTEX_SHADER)
	if err != nil {
		return err
	}
	defer gl.DeleteShader(vertShader)

	fragShader, err := compileShader(shp.FragmentShaderFilePath, gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}
	defer gl.DeleteShader(fragShader)

	programObj := gl.CreateProgram()
	gl.AttachShader(programObj, vertShader)
	gl.AttachShader(programObj, fragShader)
	gl.LinkProgram(programObj)

	var linkStatus int32
	gl.GetProgramiv(programObj, gl.LINK_STATUS, &linkStatus)
	if linkStatus == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programObj, gl.INFO_LOG_LENGTH, &logLength)
		logMsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(programObj, logLength, nil, gl.Str(logMsg))
		gl.DeleteProgram(programObj)
		return fmt.Errorf("failed to link program: %s", logMsg)
	}

	shp.programID = programObj

	return nil
}

func compileShader(file_path string, shader_type uint32) (uint32, error) {

	//load source file content as string
	source, err := os.ReadFile(file_path)
	if err != nil {
		return 0, err
	}

	//convert to c-style strings
	csource, free := gl.Strs(string(source) + "\x00") // null-terminate string

	//create new shader object
	shaderObj := gl.CreateShader(shader_type)
	gl.ShaderSource(shaderObj, 1, csource, nil)

	//free-up C memory
	free()

	gl.CompileShader(shaderObj)

	var status int32
	gl.GetShaderiv(shaderObj, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderObj, gl.INFO_LOG_LENGTH, &logLength)
		logMsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderObj, logLength, nil, gl.Str(logMsg))
		gl.DeleteShader(shaderObj)
		return 0, fmt.Errorf("failed to compile shader: %s", logMsg)
	}

	return shaderObj, nil
}
