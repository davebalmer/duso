# save_image

Save a binary image to a file.

## Syntax

```duso
save_image(binary, filename)
```

## Parameters

- `binary` - A `binary` value containing image data
- `filename` - Path where the image should be saved

## Returns

`nil`

## Description

Saves a binary image to a file. This is a convenience wrapper around `save_binary()` that exists for API symmetry with `load_image()`.

## Path Resolution

Same as `save()`. See [Files, Modules, and Paths](../files-and-modules.md#path-roots) for the full table of path roots.

## Examples

### Save a scaled image

```duso
img = load_image("original.png")
scaled = scale_image(img, 400, 400, "fill")
save_image(scaled, "thumbnail.png")
```

### Convert and save

```duso
img = load_image("input.png")
jpeg_version = convert_image(img, "jpeg")
save_image(jpeg_version, "output.jpg")
```

### Crop and save

```duso
img = load_image("photo.jpg")
quarter = crop_image(img, 0, 0, img.width / 2, img.height / 2)
save_image(quarter, "quarter.jpg")
```

## See Also

- [save_binary()](save_binary.md) - Save raw binary data
- [load_image()](load_image.md) - Load images with metadata
- [scale_image()](scale_image.md) - Resize images
- [crop_image()](crop_image.md) - Extract image regions
- [convert_image()](convert_image.md) - Convert between formats
