# load_image()

Load an image file and populate image metadata (width, height, format).

`load_image(filename)`

## Parameters

- `filename` - Path to the image file to load (PNG, JPEG, or GIF)

## Returns

A `binary` value containing the image data with populated metadata:
- `width` - Image width in pixels
- `height` - Image height in pixels
- `format` - Image format (`"png"`, `"jpeg"`, or `"gif"`)
- `content_type` - MIME type (`"image/png"`, `"image/jpeg"`, or `"image/gif"`)
- `filename` - Original filename

## Description

Loads an image file and automatically decodes it to extract width and height information. The metadata is immediately available on the returned binary value.

This is a convenience wrapper around `load_binary()` that adds image metadata extraction.

## Path Resolution

Same as `load()`. See [Files, Modules, and Paths](../files-and-modules.md#path-roots) for the full table of path roots.

## Format Support

Supports PNG, JPEG, and GIF formats. The format is automatically detected and populated in metadata.

## Examples

### Load and inspect image dimensions

```duso
img = load_image("photo.png")
print("Size:", img.width, "x", img.height)
print("Format:", img.format)
```

### Load and scale image

```duso
img = load_image("original.jpg")
scaled = scale_image(img, 800, 600, "fit")
save_image(scaled, "scaled.jpg")
```

### Load and crop

```duso
img = load_image("landscape.png")
cropped = crop_image(img, 100, 100, 400, 400)
save_image(cropped, "cropped.png")
```

## See Also

- [load_binary()](/docs/reference/load_binary.md) - Load raw binary files
- [save_image()](/docs/reference/save_image.md) - Save images to files
- [scale_image()](/docs/reference/scale_image.md) - Resize images
- [crop_image()](/docs/reference/crop_image.md) - Extract image regions
