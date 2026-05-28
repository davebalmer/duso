# convert_image

Convert an image between different formats (PNG, JPEG, GIF).

## Syntax

```duso
convert_image(image, format)
```

## Parameters

- `image` - A `binary` value containing image data (PNG, JPEG, or GIF)
- `format` - Target format: `"png"`, `"jpeg"` (or `"jpg"`), or `"gif"`

## Returns

A new `binary` value containing the image in the target format with updated metadata (format, content_type).

## Description

Converts an image from one format to another. Supports conversion between PNG, JPEG, and GIF formats. The image dimensions and content are preserved; only the encoding format changes.

## Format Support

- **PNG** - Lossless, supports full color and transparency
- **JPEG** - Lossy compression, suitable for photographs (quality: 85%)
- **GIF** - Lossless, supports animation frames (single frame images)

## Examples

### Convert PNG to JPEG

```duso
png_image = load_binary("photo.png")
jpeg = convert_image(png_image, "jpeg")
save_binary(jpeg, "photo.jpg")
```

### Convert JPEG to PNG

```duso
jpg_image = load_binary("photo.jpg")
png = convert_image(jpg_image, "png")
save_binary(png, "photo.png")
```

### Using named arguments

```duso
image = load_binary("image.gif")
converted = convert_image(image, format = "png")
```

### Format conversion in pipeline

```duso
image = load_binary("original.jpg")
scaled = scale_image(image, 800, 600, "fit")
cropped = crop_image(scaled, 100, 100, 600, 400)
result = convert_image(cropped, "png")
save_binary(result, "final.png")
```

### Convert for web delivery

```duso
original = load_binary("photo.png")
web_version = convert_image(original, "jpeg")
// JPEG typically smaller for web photos
save_binary(web_version, "web.jpg")
```

## Metadata

The returned binary includes updated metadata:

- `format` - New image format ("png", "jpeg", or "gif")
- `content_type` - MIME type ("image/png", "image/jpeg", or "image/gif")
- `width` - Image width in pixels (preserved from input)
- `height` - Image height in pixels (preserved from input)
- `filename` - Preserved from input if present

## Format Details

### PNG
- Lossless compression - no quality loss
- Supports transparency and alpha channel
- Larger file sizes for photographs

### JPEG
- Lossy compression - some quality loss
- Optimized for photographs
- Smaller file sizes, quality set to 85%
- Transparency not supported

### GIF
- Lossless compression
- Historical format, limited color palette
- Single frame images work fine
- Larger than PNG/JPEG for most content

## Performance Notes

- Memory-efficient: creates new binary only for result
- JPEG conversion uses optimized 85% quality setting for balance between size and quality
- Suitable for web workflows: format optimization, compatibility handling

## See Also

- [scale_image() - Resize images](/docs/reference/scale_image.md)
- [crop_image() - Extract image regions](/docs/reference/crop_image.md)
- [load_binary() - Load image files](/docs/reference/load_binary.md)
- [save_binary() - Save images to files](/docs/reference/save_binary.md)
- [binary - Binary data type overview](/docs/reference/binary.md)
