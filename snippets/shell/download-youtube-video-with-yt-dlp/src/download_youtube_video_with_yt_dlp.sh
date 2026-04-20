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
