package main

import (
	"fmt"
	"math"
	"runtime"
	"strings"

	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

const vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 aPos;

void main() {
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}
` + "\x00"

const fragmentShaderSource = `
#version 330 core

out vec4 FragColor;

uniform vec4 ourColor;

void main() {
	FragColor = ourColor;
}
` + "\x00"

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()


	glfw.WindowHint(glfw.ContextVersionMajor,3)
	glfw.WindowHint(glfw.ContextVersionMinor,3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	

	if err := gl.Init(); err != nil {
		print("hello")
		log.Fatalln(err)
	}

	gl.Viewport(0,0,640,480)

	glfw.GetCurrentContext().SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0,0,int32(width),int32(height))
	})

	shaderProgram, err := newShaderProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		log.Fatalln("Failed to create shader program:", err)
	}
	defer gl.DeleteProgram(shaderProgram)


	vertices := []float32{
	    -0.5, -0.5, 0.0,  // Bottom-left
	     0.5, -0.5, 0.0,  // Bottom-right
	    -0.5,  0.5, 0.0,  // Top-left
	    // -0.5,  0.5, 0.0,  // Top-left
	    //  0.5, -0.5, 0.0,  // Bottom-right
	     0.5,  0.5, 0.0,  // Top-right
	}

	indices := []uint32{
		0,1,2,
		2,3,1,
	}

	// creating vao
	var vao uint32
	gl.GenVertexArrays(1,&vao)
	gl.BindVertexArray(vao)

	// creating a VBO
	var vbo uint32;
	gl.GenBuffers(1,&vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER,len(vertices)*4,gl.Ptr(vertices),gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)


	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	// element bufer obj
	var ebo uint32
	gl.GenBuffers(1,&ebo)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER,ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER,len(indices)*4,gl.Ptr(indices),gl.STATIC_DRAW)


	// test
	var max_vertex_attr int32
	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS,&max_vertex_attr)

	fmt.Printf("max attr is %u",max_vertex_attr)

	gl.UseProgram(shaderProgram)

	ourColorUniformLocation := gl.GetUniformLocation(shaderProgram, gl.Str("ourColor"+"\x00"))
	
	// render loop
	for !window.ShouldClose() {
		
		// input
		processInput(window);

		// rendering commands
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// default: gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		// 
		
		timeValue := glfw.GetTime();
		greenValue := (math.Sin(timeValue) / 2.0) + 0.5;
		gl.Uniform4f(ourColorUniformLocation, 0.0, float32(greenValue), 0.0, 1.0)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
		
		// check and call events and swap buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

// Helper to compile shaders and link them into a program
func newShaderProgram(vertSrc, fragSrc string) (uint32, error) {
	vertShader, err := compileShader(vertSrc, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertShader)

	fragShader, err := compileShader(fragSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)
	gl.LinkProgram(program)


	return program, nil
}


func compileShader(source string, shader_type uint32) (uint32, error){

	shader := gl.CreateShader(shader_type)
	csource, free := gl.Strs(source)
	gl.ShaderSource(shader,1,csource,nil)
	free()
	gl.CompileShader(shader)

	// error handling
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		logMsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logMsg))
		return 0, fmt.Errorf("failed to compile shader: %s", logMsg)
	}

	return shader, nil
}

func processInput(window *glfw.Window){
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}