//          Copyright 2020, Vitali Baumtrok.
// Distributed under the Boost Software License, Version 1.0.
//     (See accompanying file LICENSE or copy at
//        http://www.boost.org/LICENSE_1_0.txt)

// +build texture

package main

/*
#include <string.h>

char *vertex_shader = "#version 130\n\nin vec3 vertexPosition;\nin vec2 textureCoordsIn;\nout vec2 textureCoords2;\n\nvoid main() {\n\tgl_Position = vec4(vertexPosition, 1.0f);\n\ttextureCoords2 = textureCoordsIn;\n}";
char *fragment_shader = "#version 130\n\nin vec2 textureCoords2;\nout vec4 color;\nuniform sampler2D imageTexture;\n\nvoid main() {\n\tcolor = texture(imageTexture, textureCoords2);\n}";
char *position_attribute = "vertexPosition";
char *texture_attribute = "textureCoordsIn";
char *sampler = "imageTexture";
*/
import "C"
import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
	"unsafe"
)

var (
	texVertexShader      = (**uint8)(unsafe.Pointer(&C.vertex_shader))
	texFragmentShader    = (**uint8)(unsafe.Pointer(&C.fragment_shader))
	texPositionAttribute = (*uint8)(unsafe.Pointer(C.position_attribute))
	texTextureAttribute  = (*uint8)(unsafe.Pointer(C.texture_attribute))
	texSampler           = (*uint8)(unsafe.Pointer(C.sampler))
)

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()

	if err == nil {
		var window *glfw.Window
		defer glfw.Terminate()
		window, err = glfw.CreateWindow(300, 300, "OpenGL Example", nil, nil)

		if err == nil {
			defer window.Destroy()
			window.SetKeyCallback(onKey)
			window.SetSizeCallback(onResize)
			window.MakeContextCurrent()
			err = gl.Init()

			if err == nil {
				var vTextureShader uint32
				vTextureShader, err = newShader(gl.VERTEX_SHADER, texVertexShader)

				if err == nil {
					var fTextureShader uint32
					fTextureShader, err = newShader(gl.FRAGMENT_SHADER, texFragmentShader)

					if err == nil {
						var textureProgram uint32
						textureProgram, err = newProgram(vTextureShader, fTextureShader)

						if err == nil {
							defer gl.DeleteProgram(textureProgram)
							vbos := newVBOs(2)
							defer gl.DeleteBuffers(int32(len(vbos)), &vbos[0])
							vaos := newVAOs(1)
							defer gl.DeleteVertexArrays(int32(len(vaos)), &vaos[0])
							textures := newTextures(1)
							defer gl.DeleteTextures(int32(len(textures)), &textures[0])

							bindTextureObjects(textureProgram, vaos, vbos, textures)
							gl.UseProgram(textureProgram)

							//textureSampler := int32(gl.GetAttribLocation(textureProgram, texSampler))

							// transparency
							// gl.Enable(gl.BLEND);
							// gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);

							// wireframe mode
							// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

							for !window.ShouldClose() {
								gl.ClearColor(0, 0, 0, 0)
								gl.Clear(gl.COLOR_BUFFER_BIT)
								gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
								/*
									gl.UseProgram(textureProgram)
									gl.ActiveTexture(gl.TEXTURE0);
									gl.BindTexture(gl.TEXTURE_2D, textures[0])
									gl.Uniform1i(int32(textureSampler), 0)
									gl.BindVertexArray(vaos[0])

									gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, unsafe.Pointer(nil))
								*/

								window.SwapBuffers()
								glfw.PollEvents()
							}
						}
					}
				}
			}
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

func onKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func onResize(w *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func newShader(shaderType uint32, shaderSource **uint8) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, shaderSource, nil)
	gl.CompileShader(shader)
	err := checkShader(shader, gl.COMPILE_STATUS)

	if err != nil {
		gl.DeleteShader(shader)
	}
	return shader, err
}

func newProgram(vShader, fShader uint32) (uint32, error) {
	program := gl.CreateProgram()
	gl.AttachShader(program, vShader)
	gl.AttachShader(program, fShader)
	gl.LinkProgram(program)
	err := checkProgram(program, gl.LINK_STATUS)

	if err == nil {
		gl.ValidateProgram(program)
		err = checkProgram(program, gl.VALIDATE_STATUS)

		if err == nil {
			gl.EnableVertexAttribArray(0)
			gl.EnableVertexAttribArray(1)

		} else {
			gl.DeleteProgram(program)
		}
	}
	return program, err
}

func checkShader(shader, statusType uint32) error {
	var status int32
	var err error

	gl.GetShaderiv(shader, statusType, &status)

	if status == gl.FALSE {
		var length int32
		var infoLog string

		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)

		if length > 0 {
			infoLogBytes := make([]byte, length)
			gl.GetShaderInfoLog(shader, length, nil, &infoLogBytes[0])
			infoLog = string(infoLogBytes)
		}
		err = errors.New("shader " + infoLog)
	}
	return err
}

func checkProgram(program, statusType uint32) error {
	var status int32
	var err error

	gl.GetProgramiv(program, statusType, &status)

	if status == gl.FALSE {
		var length int32
		var infoLog string

		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)

		if length > 0 {
			infoLogBytes := make([]byte, length)
			gl.GetProgramInfoLog(program, length, nil, &infoLogBytes[0])
			infoLog = string(infoLogBytes)
		}
		err = errors.New("program " + infoLog)
	}
	return err
}

func newVBOs(n int) []uint32 {
	vbos := make([]uint32, n)
	gl.GenBuffers(int32(len(vbos)), &vbos[0])
	return vbos
}

func newVAOs(n int) []uint32 {
	vaos := make([]uint32, n)
	gl.GenVertexArrays(int32(len(vaos)), &vaos[0])
	return vaos
}

func newTextures(n int) []uint32 {
	textures := make([]uint32, n)
	gl.GenTextures(int32(len(textures)), &textures[0])
	return textures
}

func bindTextureObjects(program uint32, vaos, vbos, textures []uint32) {
	textureData := newTextureData()
	positionLocation := uint32(gl.GetAttribLocation(program, texPositionAttribute))
	textureLocation := uint32(gl.GetAttribLocation(program, texTextureAttribute))
	// x, y, z, x_tex, y_tex (two triangles)
	vertices := []float32{
		0.5, 0.5, 0.0, 1.0, 1.0,
		0.5, 0.0, 0.0, 1.0, 0.0,
		0.0, 0.5, 0.0, 0.0, 1.0,
		0.0, 0.0, 0.0, 0.0, 0.0,
	}
	gl.BindTexture(gl.TEXTURE_2D, textures[0])
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 64, 64, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(textureData))
	gl.BindBuffer(gl.ARRAY_BUFFER, vbos[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(positionLocation)
	gl.EnableVertexAttribArray(textureLocation)
	gl.VertexAttribPointer(positionLocation, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(textureLocation, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	/*
	textureData := newTextureData()
	positionLocation := uint32(gl.GetAttribLocation(program, texPositionAttribute))
	textureLocation := uint32(gl.GetAttribLocation(program, texTextureAttribute))
	// x, y, z, x_tex, y_tex (two triangles)
	vertices := []float32{
		0.5, 0.5, 0.0, 1.0, 1.0,
		0.5, 0.0, 0.0, 1.0, 0.0,
		0.0, 0.5, 0.0, 0.0, 1.0,
		0.0, 0.0, 0.0, 0.0, 0.0,
	}
	indices := []uint32 {
		0, 1, 2,
		2, 1, 3,
	}

	gl.BindVertexArray(vaos[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, vbos[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vbos[1])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// position
	gl.VertexAttribPointer(positionLocation, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(positionLocation)
	// color
	gl.VertexAttribPointer(textureLocation, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(textureLocation)

	gl.ActiveTexture(gl.TEXTURE0);
	gl.BindTexture(gl.TEXTURE_2D, textures[0]);

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 64, 64, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(textureData));
	// gl.GenerateMipmap(gl.TEXTURE_2D);
	*/
}

func newTextureData() []uint8 {
	data := make([]uint8, 64*64*4)
	for i := 0; i < 64*64; i++ {
		offset := i * 4
		if (i/16+i/(16*64))%2 == 0 {
			data[offset] = 255
			data[offset+1] = 255
			data[offset+2] = 255
			data[offset+3] = 255
		}
	}
	return data
}
