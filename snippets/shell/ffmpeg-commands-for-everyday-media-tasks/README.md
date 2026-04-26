Save time on common video and audio work with `ffmpeg` commands for conversion, trimming, resizing, screenshots, GIF previews, audio extraction, compression, and subtitles.

## What This Snippet Covers

- Converting a video to MP4
- Extracting audio to MP3
- Trimming a clip by time range
- Resizing video while keeping aspect ratio
- Making a small GIF preview
- Capturing a still frame
- Compressing a larger video into H.264 MP4
- Burning subtitles into a video file

## Before Using

- Replace the example file names with media files that exist locally.
- Run output-producing commands from a writable directory.
- Keep a backup copy of source files before testing destructive overwrite patterns.

## Code

```sh
ffmpeg -i input.mov output.mp4
# Convert a video to MP4 with default codec choices.

ffmpeg -i input.mp4 -vn output.mp3
# Extract audio only and save it as MP3.

ffmpeg -ss 00:00:10 -to 00:00:25 -i input.mp4 -c copy clip.mp4
# Trim a clip without re-encoding when the source format allows it.

ffmpeg -i input.mp4 -vf scale=1280:-2 resized.mp4
# Resize video width to 1280 while keeping the aspect ratio.

ffmpeg -i input.mp4 -vf 'fps=12,scale=720:-2:flags=lanczos' preview.gif
# Convert a short clip into a GIF preview.

ffmpeg -i input.mp4 -vframes 1 cover.jpg
# Capture the first frame as an image.

ffmpeg -i input.mov -c:v libx264 -crf 23 -preset medium -c:a aac compressed.mp4
# Re-encode to a smaller H.264 MP4 file.

ffmpeg -i input.mp4 -vf subtitles=subtitles.srt subtitled.mp4
# Burn subtitles directly into the video.
```

## Why These Commands Are Useful

- They cover the media tasks people repeatedly search for in docs or old notes.
- They keep one useful command shape per task instead of burying the reader in codec theory.
- They provide a practical starting point for local editing, exports, and previews.

## Limitations

- This snippet stays `Draft` because it depends on placeholder media files and a local `ffmpeg` install.
- `-c copy` trimming works best around compatible keyframe boundaries.
- Subtitle filters and available codecs can vary by build and operating system.

## Manual Verification

1. Confirm `ffmpeg -version` works.
2. Replace the example file names with real local media files.
3. Run each command from a writable directory.
4. Confirm the generated media output matches the intended task.

## Files

- `src/ffmpeg_commands_for_everyday_media_tasks.sh`
- `snippet.json`