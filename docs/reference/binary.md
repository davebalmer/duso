# binary

Immutable binary data type for handling files, images, and other binary content.

## Overview

A `binary` value represents immutable binary data (raw bytes) loaded from files or created within scripts. Binary values are memory-efficient - they use pointer sharing when passed between script scopes, so multiple workers can hold references to the same binary data without copying.

## Working with Binary Data

- **Load files** - [load_binary()](/docs/reference/load_binary.md)
- **Save files** - [save_binary()](/docs/reference/save_binary.md)
- **Text encoding** - [encode_base64()](/docs/reference/encode_base64.md) and [decode_base64()](/docs/reference/decode_base64.md)

## Common Operations

### Get size in bytes

```duso
image = load_binary("photo.png")
print(len(image))  // prints byte count
```

### Access metadata

```duso
file = load_binary("photo.png")
filename = file["filename"]  // original filename
```

**Available metadata:**
- `filename` - original filename from `load_binary()` or HTTP file uploads
- `content_type` - MIME type, may be populated by HTTP file uploads

## HTTP File Uploads

When file uploads are enabled, uploaded files are accessible via `req.files`:

```duso
ctx = context()
req = ctx.request()

// Single file upload
file = req.files.avatar
if file then
  print(file.filename)      // "avatar.png"
  print(file.content_type)  // "image/png"
  print(file.size)          // bytes
  
  if type(file.data) == "binary" then
    save_binary(file.data, "/STORE/uploads/" + file.filename)
  elseif type(file.data) == "string" then
    // Text file (JSON, XML, etc)
    parsed = parse_json(file.data)
  end
end

// Multiple files on same field → array
for f in req.files.attachments do
  print(f.filename)
end
```

**MIME Type Handling:**
- Text MIME types (`text/*`, `application/json`, `application/xml`, etc.) → `file.data` is a string
- Binary MIME types (images, archives, etc.) → `file.data` is a binary value

**File object properties:**
- `data` - File content (binary or string depending on MIME type)
- `filename` - Original filename
- `content_type` - MIME type
- `size` - File size in bytes

## Type Checking

```duso
data = load_binary("file.bin")
if type(data) == "binary" then
  print("It's binary!")
end
```

## Truthiness

```duso
image = load_binary("photo.png")
if image then
  print("Has data")
end
```

## Properties

- **Immutable** - Cannot be modified after creation
- **No indexing** - Cannot access individual bytes like arrays
- **Pointer-based** - Multiple workers can reference the same data without copying
- **Memory efficient** - Only the pointer is shared, not the raw bytes
- **GC managed** - Automatically cleaned up when last reference is dropped

## Limitations

- **No partial reads** - Must load entire file into memory
- **No modification** - Binary values are immutable
- **No direct inspection** - Cannot index individual bytes; use metadata instead
- **Text encoding** - For text files, use [load()](/docs/reference/load.md) instead for UTF-8 string handling

## See Also

- [load_binary() - Load binary files](/docs/reference/load_binary.md)
- [save_binary() - Save binary files](/docs/reference/save_binary.md)
- [encode_base64() - Encode binary to base64 text](/docs/reference/encode_base64.md)
- [decode_base64() - Decode base64 text to binary](/docs/reference/decode_base64.md)
- [load() / save() - Text file operations](/docs/reference/load.md)
- [type() - Type checking](/docs/reference/type.md)
- [len() - Length operations](/docs/reference/len.md)
- [http_server() - HTTP file uploads](/docs/reference/http_server.md)
