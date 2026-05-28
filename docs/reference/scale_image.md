# scale_image

Scale an image to a target size with different aspect ratio handling modes.

## Syntax

```duso
scale_image(image, max_x, max_y, mode)
```

## Parameters

- `image` - A `binary` value containing image data (PNG, JPEG, or GIF)
- `max_x` - Target width as a number (must be positive)
- `max_y` - Target height as a number (must be positive)
- `mode` - Scaling mode: `"fit"`, `"fill"`, or `"stretch"`

## Returns

A new `binary` value containing the scaled image with updated metadata (width, height, format, content_type).

## Description

Scales an image to fit within or fill the specified dimensions. The mode parameter controls how aspect ratio is handled:

- **`"fit"`** (default) - Scale to fit within max_x × max_y while preserving aspect ratio. Result may be smaller than target dimensions.
- **`"fill"`** - Scale to fill max_x × max_y while preserving aspect ratio. Result is cropped to exact dimensions from center.
- **`"stretch"`** - Scale to exact max_x × max_y, ignoring aspect ratio. May distort image.

## Format Support

Supports PNG, JPEG, and GIF formats. Output format matches input format.

## Examples

### Scale for thumbnail (preserve aspect)

```duso
image = load_binary("photo.jpg")
thumbnail = scale_image(image, 256, 256, "fit")
save_binary(thumbnail, "photo_thumb.jpg")
```

### Scale for avatar (fill and crop)

```duso
avatar = load_binary("avatar.png")
square = scale_image(avatar, 200, 200, "fill")
save_binary(square, "avatar_200x200.png")
```

### Stretch to exact dimensions

```duso
image = load_binary("image.gif")
stretched = scale_image(image, 800, 600, "stretch")
```

### Using named arguments

```duso
image = load_binary("photo.jpg")
scaled = scale_image(image, 512, 512, mode = "fit")
```

### Pipelining with other operations

```duso
image = load_binary("photo.jpg")
scaled = scale_image(image, 400, 400, "fill")
converted = convert_image(scaled, "png")
save_binary(converted, "output.png")
```

## Metadata

The returned binary includes updated metadata:

- `width` - New image width in pixels
- `height` - New image height in pixels
- `format` - Image format ("png", "jpeg", or "gif")
- `content_type` - MIME type ("image/png", "image/jpeg", or "image/gif")
- `filename` - Preserved from input if present

## Performance Notes

- Uses nearest-neighbor sampling for efficiency
- Memory-efficient: creates new binary only for result, original is garbage-collected when unused
- Suitable for web workflows: avatar processing, thumbnail generation, image resizing

## See Also

- [crop_image() - Extract rectangular regions](/docs/reference/crop_image.md)
- [convert_image() - Convert between image formats](/docs/reference/convert_image.md)
- [load_binary() - Load image files](/docs/reference/load_binary.md)
- [save_binary() - Save images to files](/docs/reference/save_binary.md)
- [binary - Binary data type overview](/docs/reference/binary.md)
