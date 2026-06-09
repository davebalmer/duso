# crop_image()

Extract a rectangular region from an image.

`crop_image(image, x, y, width, height)`

## Parameters

- `image` - A `binary` value containing image data (PNG, JPEG, or GIF)
- `x` - X coordinate of the top-left corner (0-based)
- `y` - Y coordinate of the top-left corner (0-based)
- `width` - Width of the region to extract (must be positive)
- `height` - Height of the region to extract (must be positive)

## Returns

A new `binary` value containing the cropped image with updated metadata (width, height, format, content_type).

## Description

Extracts a rectangular region from an image, starting at coordinates (x, y) with the specified dimensions. The crop region is automatically clipped to the image bounds.

## Format Support

Supports PNG, JPEG, and GIF formats. Output format matches input format.

## Examples

### Extract center region

```duso
image = load_image("photo.jpg")
center = crop_image(image, 100, 100, 200, 200)
save_image(center, "center.jpg")
```

### Create square crop

```duso
portrait = load_image("portrait.png")
square = crop_image(portrait, 0, 50, 300, 300)
save_image(square, "square.png")
```

### Using named arguments

```duso
image = load_image("photo.gif")
cropped = crop_image(image, x = 50, y = 50, width = 150, height = 150)
```

### Positional then named

```duso
image = load_image("image.jpg")
region = crop_image(image, 10, 10, width = 100, height = 100)
```

### Pipelining

```duso
image = load_image("photo.jpg")
cropped = crop_image(image, 50, 50, 400, 400)
thumbnail = scale_image(cropped, 200, 200, "fit")
save_image(thumbnail, "thumb.jpg")
```

## Behavior

- **Bounds clipping** - If the specified region extends beyond image boundaries, it's automatically clipped to fit
- **Zero-based indexing** - Top-left corner is (0, 0)
- **Exact dimensions** - Result is exactly width × height pixels (or smaller if clipped)

## Metadata

The returned binary includes updated metadata:

- `width` - New image width in pixels (result of crop)
- `height` - New image height in pixels (result of crop)
- `format` - Image format ("png", "jpeg", or "gif")
- `content_type` - MIME type ("image/png", "image/jpeg", or "image/gif")
- `filename` - Preserved from input if present

## Performance Notes

- Memory-efficient: creates new binary only for result, original is garbage-collected when unused
- Suitable for web workflows: avatar cropping, region extraction, image composition

## See Also

- [scale_image() - Resize images](/docs/reference/scale_image.md)
- [convert_image() - Convert between image formats](/docs/reference/convert_image.md)
- [load_image() - Load image files](/docs/reference/load_image.md)
- [save_image() - Save images to files](/docs/reference/save_image.md)
- [binary - Binary data type overview](/docs/reference/binary.md)
