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