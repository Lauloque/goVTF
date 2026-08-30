# goVTF

## What is this?

goVTF is a humble attempt to imitate VTFLib and VTFCmd, Valve's library and CLI for creating Valve Texture Format (VTF) files, with Linux support.

It was created in the context of [Boot.dev's Interactive Course: First Personal Project](https://www.boot.dev/courses/build-personal-project-1).

<img width="900" height="550" alt="image" src="https://github.com/user-attachments/assets/606dbef9-033b-428b-9f5c-3f3e41870c3a" />

The current implementation converts PNG and JPEG images to VTF 7.3. It supports DXT1 output and DXT5 output for images requiring transparency.

## Current usage

Long options use two hyphens. Although the CLI also provides the conventional
short forms `-h` and `-o`, the examples below consistently use long options.

```text
goVTF file <path> [options]
```

### Implemented options

- `--output <directory>` — output directory (default: `./`).
- `--alphaformat <format>` — high-resolution output format. Supported values
  are `dxt1` and `dxt5` (default: `dxt1`). Use `dxt5` to retain an image's
  full alpha channel.
- `--help` — display help.

The input path is a positional argument to the `file` command; `--file` is not
an implemented option.

### Examples

Convert an opaque image using the default DXT1 format:

```bash
goVTF file image.png --output ./
```

Convert a PNG while retaining full transparency with DXT5:

```bash
goVTF file image.png --alphaformat dxt5 --output ./
```

Use `goVTF --help` or `goVTF file --help` for generated command help.

## Potential future options

The following options are inspired by VTFCmd. They are ideas for future work
and are **not currently implemented or accepted by goVTF**:

```text
--folder <path>             Input directory search string.
--prefix <string>           Output file prefix.
--postfix <string>          Output file postfix.
--version <string>          Output VTF version.
--format <string>           Output format for non-alpha textures.
--flag <string>             Output flag to set.
--resize                    Resize the input to a power of two.
--rmethod <string>          Resize method.
--rfilter <string>          Resize filter.
--rsharpen <string>         Resize sharpening filter.
--rwidth <integer>          Resize to a specific width.
--rheight <integer>         Resize to a specific height.
--rclampwidth <integer>     Maximum resize width.
--rclampheight <integer>    Maximum resize height.
--gamma                     Gamma-correct the image.
--gcorrection <number>      Gamma correction value.
--nomipmaps                 Do not generate mipmaps.
--mfilter <string>          Mipmap filter.
--msharpen <string>         Mipmap sharpening filter.
--normal                    Convert the input to a normal map.
--nkernel <string>          Normal-map generation kernel.
--nheight <string>          Normal-map height calculation.
--nalpha <string>           Normal-map alpha result.
--nscale <number>           Normal-map scale.
--nwrap                     Wrap the normal map for tiled textures.
--bumpscale <number>        Engine bump-mapping scale.
--nothumbnail               Do not generate a thumbnail image.
--noreflectivity            Do not calculate reflectivity.
--shader <string>           Create a material for the texture.
--param <string> <string>   Add a parameter to the material.
--recurse                   Process directories recursively.
--exportformat <string>     Export VTF files to another image format.
--silent                    Silent mode.
--pause                     Pause when finished.
```

## Resources

- [VTF (Valve Texture Format) - Valve Developer Community](<https://developer.valvesoftware.com/w/index.php?title=VTF_(Valve_Texture_Format)&utm_source=chatgpt.com>)

- [VMT - Valve Developer Community](https://developer.valvesoftware.com/wiki/VMT)

- [VTFCmd - Valve Developer Community](https://developer.valvesoftware.com/wiki/VTFCmd)

- [VTFLib - Valve Developer Community](https://developer.valvesoftware.com/wiki/VTFLib)

- [Sky-rym/VTFEdit-Reloaded](https://github.com/Sky-rym/VTFEdit-Reloaded)

- [VTF Spray Converter](https://rafradek.github.io/Mishcatt/)

### License

Project: see [LICENSE](LICENSE).

VTF.jpeg: test file taken from [Steam 社群 :: :: VTF Edit-Chan](https://steamcommunity.com/sharedfiles/filedetails/?&id=3046622181)
