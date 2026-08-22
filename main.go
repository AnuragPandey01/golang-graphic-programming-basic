package main

import (
	"fmt"
	"runtime"

	"log"

	"github.com/anurag01/gogine/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

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

	shaderProgram := shader.ShaderProgram{
		FragmentShaderFilePath: "assets/fragment_shader.glsl",
		VertexShaderFilePath: "assets/vertex_shader.glsl",
	}

	err = shaderProgram.Compile()
	if err != nil {
		log.Fatalln("Failed to create shader program:", err)
	}
	defer shaderProgram.Delete()


	vertices := []float32{
		// positions        // color
	    -0.5, -0.5, 0.0, 	1.0, 0.0, 0.0,  	// Bottom-left
	     0.5, -0.5, 0.0, 	0.0, 1.0, 0.0, 		// Bottom-right
	    -0.5,  0.5, 0.0, 	0.0, 0.0, 1.0,  	// Top-left
	    // -0.5,  0.5, 0.0,  					// Top-left
	    //  0.5, -0.5, 0.0,  					// Bottom-right
	     0.5,  0.5, 0.0, 	1.0, 0.0, 0.0,   	// Top-right
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

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4*2, 0)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 3*4*2, 3*4)
	gl.EnableVertexAttribArray(1)


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

	shaderProgram.Use()

	
	// render loop
	for !window.ShouldClose() {
		
		// input
		processInput(window);

		// rendering commands
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// default: gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
		
		// check and call events and swap buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func processInput(window *glfw.Window){
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}