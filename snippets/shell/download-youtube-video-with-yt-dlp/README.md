# Download YouTube Video with yt-dlp

A commented `yt-dlp` command for downloading one YouTube video with resumable fragments, explicit format sorting, and optional proxy support.

## What This Snippet Covers

- Verbose logs for extractor and network debugging
- Explicit resume behavior for partial downloads
- `--no-playlist` protection for watch URLs that also contain a playlist id
- Separate best video plus best audio with a combined-format fallback
- Parallel fragment downloads for DASH and HLS streams
- Stable id-based output filenames
- A predictable merged container after `ffmpeg` combines streams

## Before Using

- Replace `http://your-proxy.example:8080` with your real proxy URL, or remove the `--proxy` line entirely.
- Replace `https://www.youtube.com/watch?v=VIDEO_ID` with the target video URL.
- Install `ffmpeg` because `bestvideo+bestaudio` merging depends on it.

## Code

```sh
yt-dlp -v --continue --no-playlist \
  --format "bestvideo+bestaudio/b" \
  --concurrent-fragments 5 \
  --format-sort "quality,res,fps,tbr,vbr,size" \
  --merge-output-format mkv \
  -o "%(id)s.%(ext)s" \
  --proxy "http://your-proxy.example:8080" \
  "https://www.youtube.com/watch?v=VIDEO_ID"

# --continue: Keep resume behavior explicit for partial downloads.
# --no-playlist: Avoid downloading a full playlist when the URL also contains list=.
# --format "bestvideo+bestaudio/b": Prefer best separate video and audio, then fall back to the best combined format.
# --concurrent-fragments 5: Use moderate parallelism for DASH or HLS fragment downloads.
# --format-sort "quality,res,fps,tbr,vbr,size": Prefer quality first, then resolution, frame rate, bitrate, and size.
# --merge-output-format mkv: Keep the final merged container predictable.
# -o "%(id)s.%(ext)s": Use short, stable filenames based on the video id.
# --proxy: Replace this value or remove the line if direct network access already works.
```

## Why These Extra Options Were Added

- `--no-playlist` is useful for YouTube watch URLs that carry an extra `list=` query parameter, because it prevents accidental playlist downloads.
- `--format-sort "quality,res,fps,tbr,vbr,size"` keeps video quality ahead of raw file size while still using size as a later tie-breaker.
- `--merge-output-format mkv` makes the final container more predictable when `yt-dlp` downloads separate video and audio streams.
- The proxy line stays in the snippet because your original use case included it, but it is intentionally easy to remove.

## Limitations

- This snippet was not executed in the repository environment, so it stays `Draft`.
- yt-dlp behavior can change when YouTube changes extractor behavior, available formats, or throttling rules.
- If your proxy requires authentication, use a full proxy URL such as `http://user:pass@proxy.example:8080`.
- If `ffmpeg` is missing, `bestvideo+bestaudio` downloads may not merge into one final media file.

## Manual Verification

1. Confirm `yt-dlp --version` works.
2. Confirm `ffmpeg -version` works.
3. Replace the proxy and video URL placeholders, or remove the proxy line.
4. Run the command from a writable directory.
5. Confirm the final file is downloaded and playable.

## Files

- `src/download_youtube_video_with_yt_dlp.sh`
- `snippet.json`