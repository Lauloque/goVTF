# goVTF

## What is this?

goVTF is a humble attempt to imit ate VTFLib and VTFCmd, VALVe's library and its CLI to create a "Valve Texture File" (VTF) used in VALVe games. But usable in Linux.

It is made in the context of [Boot.dev's Interactive Course: First Personal Project ](https://www.boot.dev/courses/build-personal-project-1).

First implementation should just be able to take a common picture file and output a VTF usable as a "spray" in TF2.

## Resources

- [VTF (Valve Texture Format) - Valve Developer Community](<https://developer.valvesoftware.com/w/index.php?title=VTF_(Valve_Texture_Format)&utm_source=chatgpt.com>)

- [VMT - Valve Developer Community](https://developer.valvesoftware.com/wiki/VMT)

- [VTFCmd - Valve Developer Community](https://developer.valvesoftware.com/wiki/VTFCmd)

- [VTFLib - Valve Developer Community](https://developer.valvesoftware.com/wiki/VTFLib)

- [Sky-rym/VTFEdit-Reloaded](https://github.com/Sky-rym/VTFEdit-Reloaded)

## cmd reference:

```
-file <path>             (Input file path.)
-folder <path>           (Input directory search string.)
-output <path>           (Output directory.)
-prefix <string>         (Output file prefix.)
-postfix <string>        (Output file postfix.)
-version <string>        (Ouput version.)
-format <string>         (Ouput format to use on non-alpha textures.)
-alphaformat <string>    (Ouput format to use on alpha textures.)
-flag <string>           (Output flags to set.)
-resize                  (Resize the input to a power of 2.)
-rmethod <string>        (Resize method to use.)
-rfilter <string>        (Resize filter to use.)
-rsharpen <string>       (Resize sharpen filter to use.)
-rwidth <integer>        (Resize to specific width.)
-rheight <integer>       (Resize to specific height.)
-rclampwidth <integer>   (Maximum width to resize to.)
-rclampheight <integer>  (Maximum height to resize to.)
-gamma                   (Gamma correct image.)
-gcorrection <single>    (Gamma correction to use.)
-nomipmaps               (Don't generate mipmaps.)
-mfilter <string>        (Mipmap filter to use.)
-msharpen <string>       (Mipmap sharpen filter to use.)
-normal                  (Convert input file to normal map.)
-nkernel <string>        (Normal map generation kernel to use.)
-nheight <string>        (Normal map height calculation to use.)
-nalpha <string>         (Normal map alpha result to set.)
-nscale <single>         (Normal map scale to use.)
-nwrap                   (Wrap the normal map for tiled textures.)
-bumpscale <single>      (Engine bump mapping scale to use.)
-nothumbnail             (Don't generate thumbnail image.)
-noreflectivity          (Don't calculate reflectivity.)
-shader <string>         (Create a material for the texture.)
-param <string> <string> (Add a parameter to the material.)
-recurse                 (Process directories recursively.)
-exportformat <string>   (Convert VTF files to the format of this extension.)
-silent                  (Silent mode.)
-pause                   (Pause when done.)
-help                    (Display vtfcmd help.)
```

### License

Project: see [LICENSE]()
VTF.jpeg: test file taken from [Steam 社群 :: :: VTF Edit-Chan](https://steamcommunity.com/sharedfiles/filedetails/?&id=3046622181)
