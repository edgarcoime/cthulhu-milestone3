resource "aws_s3_bucket" "cthulhu_platform" {
  bucket = "cthulhu-platform"
}

resource "aws_s3_bucket_cors_configuration" "cthulhu_platform" {
  bucket = aws_s3_bucket.cthulhu_platform.id

  cors_rule {
    allowed_methods = ["GET", "PUT", "HEAD"]
    # TODO: Make stricter for production
    allowed_origins = ["*"]
    allowed_headers = ["*"]
    expose_headers  = ["ETag"]
  }
}
