from typing import Dict

from fastapi import FastAPI, File, HTTPException, UploadFile


MAX_FILE_SIZE_BYTES = 1024 * 1024
ALLOWED_CONTENT_TYPES = {"text/plain", "text/csv"}

app = FastAPI()


@app.post("/upload")
async def upload_file(file: UploadFile = File(...)) -> Dict[str, object]:
    """Validate one uploaded file before accepting it."""
    # Flow:
    #   validate filename and content type
    #      |
    #      +-> invalid upload -> raise HTTPException
    #      `-> read body -> validate size -> return accepted file metadata
    sanitized_filename = file.filename.replace("\\", "/").split("/")[-1].strip() if file.filename else ""
    if sanitized_filename == "":
        raise HTTPException(status_code=400, detail="filename is required")
    if file.content_type not in ALLOWED_CONTENT_TYPES:
        raise HTTPException(status_code=415, detail="unsupported content type")

    try:
        file_size = 0
        while True:
            chunk = await file.read(64 * 1024)
            if chunk == b"":
                break
            file_size += len(chunk)
            if file_size > MAX_FILE_SIZE_BYTES:
                raise HTTPException(status_code=413, detail="file is too large")
    finally:
        await file.close()

    if file_size == 0:
        raise HTTPException(status_code=400, detail="file is empty")

    return {
        "filename": sanitized_filename,
        "content_type": file.content_type,
        "size": file_size,
    }


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)