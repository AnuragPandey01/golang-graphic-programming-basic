package main

import (
	"runtime"

	"log"

	"github.com/anurag01/gogine/shader"
	"github.com/anurag01/gogine/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/go-gl/mathgl/mgl32"

	_ "image/jpeg"
	_ "image/png"
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

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
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

	fbWidth, fbHeight := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))

	glfw.GetCurrentContext().SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	shaderProgram := shader.NewShaderProgram("assets/vertex_shader.glsl","assets/fragment_shader.glsl")

	err = shaderProgram.Compile()
	if err != nil {
		log.Fatalln("Failed to create shader program:", err)
	}
	defer shaderProgram.Delete()

	vertices := []float32{
		// positions          // texture coords
		// Front
		-0.5, -0.5,  0.5,      0.0, 0.0,
		 0.5, -0.5,  0.5,      1.0, 0.0,
		 0.5,  0.5,  0.5,      1.0, 1.0,
	
		 0.5,  0.5,  0.5,      1.0, 1.0,
		-0.5,  0.5,  0.5,      0.0, 1.0,
		-0.5, -0.5,  0.5,      0.0, 0.0,
	
		// Back
		 0.5, -0.5, -0.5,      0.0, 0.0,
		-0.5, -0.5, -0.5,      1.0, 0.0,
		-0.5,  0.5, -0.5,      1.0, 1.0,
	
		-0.5,  0.5, -0.5,      1.0, 1.0,
		 0.5,  0.5, -0.5,      0.0, 1.0,
		 0.5, -0.5, -0.5,      0.0, 0.0,
	
		// Left
		-0.5, -0.5, -0.5,      0.0, 0.0,
		-0.5, -0.5,  0.5,      1.0, 0.0,
		-0.5,  0.5,  0.5,      1.0, 1.0,
	
		-0.5,  0.5,  0.5,      1.0, 1.0,
		-0.5,  0.5, -0.5,      0.0, 1.0,
		-0.5, -0.5, -0.5,      0.0, 0.0,
	
		// Right
		 0.5, -0.5,  0.5,      0.0, 0.0,
		 0.5, -0.5, -0.5,      1.0, 0.0,
		 0.5,  0.5, -0.5,      1.0, 1.0,
	
		 0.5,  0.5, -0.5,      1.0, 1.0,
		 0.5,  0.5,  0.5,      0.0, 1.0,
		 0.5, -0.5,  0.5,      0.0, 0.0,
	
		// Top
		-0.5,  0.5,  0.5,      0.0, 0.0,
		 0.5,  0.5,  0.5,      1.0, 0.0,
		 0.5,  0.5, -0.5,      1.0, 1.0,
	
		 0.5,  0.5, -0.5,      1.0, 1.0,
		-0.5,  0.5, -0.5,      0.0, 1.0,
		-0.5,  0.5,  0.5,      0.0, 0.0,
	
		// Bottom
		-0.5, -0.5, -0.5,      0.0, 0.0,
		 0.5, -0.5, -0.5,      1.0, 0.0,
		 0.5, -0.5,  0.5,      1.0, 1.0,
	
		 0.5, -0.5,  0.5,      1.0, 1.0,
		-0.5, -0.5,  0.5,      0.0, 1.0,
		-0.5, -0.5, -0.5,      0.0, 0.0,
	}

	indices := []uint32{
		0, 1, 2,
		2, 3, 1,
	}

	//creating a texture
	var texture1, texture2 uint32
	gl.GenTextures(1, &texture1)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture1)
	tex1imgByteArray, tex1bounds, err := utils.LoadImageAsset("assets/container.jpg")
	if err != nil {
		log.Fatalf("Failed loading texture asset, %v", err)
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(tex1bounds.Dx()), int32(tex1bounds.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(tex1imgByteArray))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	
	gl.GenTextures(1, &texture2)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, texture2)
	tex2imgByteArray, tex2bounds, err := utils.LoadImageAsset("assets/awesomeface.png")
	if err != nil {
		log.Fatalf("Failed loading texture asset, %v", err)
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(tex2bounds.Dx()), int32(tex2bounds.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(tex2imgByteArray))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// creating vao
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// creating a VBO
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, (3+2)*4, 0)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, (3+2)*4, 3*4)
	gl.EnableVertexAttribArray(1)

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)
	

	// element bufer obj
	var ebo uint32
	gl.GenBuffers(1, &ebo)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	shaderProgram.Use()

	gl.Uniform1i(shaderProgram.GetUniformLocation("texture1"), 0)
	gl.Uniform1i(shaderProgram.GetUniformLocation("texture2"), 1)


	view := mgl32.Ident4()
	view = view.Mul4(mgl32.Translate3D(0, 0, -3))
	gl.UniformMatrix4fv(shaderProgram.GetUniformLocation("view"), 1, false, &view[0])

	projection := mgl32.Perspective(mgl32.DegToRad(45), float32(fbWidth)/float32(fbHeight), 0.1, 100)
	gl.UniformMatrix4fv(shaderProgram.GetUniformLocation("projection"), 1, false, &projection[0])

	cubePositions := []mgl32.Vec3{
		{0, 0, 0},
		{2, 5, -15},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{2.4, -0.4, -3.5},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}
	
	// render loop
	for !window.ShouldClose() {

		// input
		processInput(window)

		// Update projection matrix in case of resize
		w, h := window.GetFramebufferSize()
		if w > 0 && h > 0 {
			proj := mgl32.Perspective(mgl32.DegToRad(45), float32(w)/float32(h), 0.1, 100)
			gl.UniformMatrix4fv(shaderProgram.GetUniformLocation("projection"), 1, false, &proj[0])
		}

		// rendering commands
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);

		// default: gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

		// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
		for i := 0; i<10; i++ {
			model := mgl32.Ident4()
			model = model.Mul4(mgl32.Translate3D(cubePositions[i][0], cubePositions[i][1], cubePositions[i][2]))
			angle := 20.00 * float32(i)
			model = model.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(angle),mgl32.Vec3{1.0, 0.3, 0.5}))
			gl.UniformMatrix4fv(shaderProgram.GetUniformLocation("model"), 1, false, &model[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		model := mgl32.Ident4()
		model = model.Mul4(mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{1.0, 0.0, 0.0}))
		gl.UniformMatrix4fv(shaderProgram.GetUniformLocation("model"), 1, false, &model[0])

		// check and call events and swap buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}
