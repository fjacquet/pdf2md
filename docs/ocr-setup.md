# OCR Setup

`pdf2md` falls back to OCR when a PDF page contains no extractable text
(scanned document, image-only page). The OCR pipeline uses PaddleOCR's
PP-OCRv5 mobile models via ONNX Runtime.

Unlike the DocLayout-YOLO model used for layout detection, PaddleOCR ONNX
builds are not available from stable community mirrors, so `pdf2md` does
**not** auto-download them. Install them once as described below.

## 1. Install ONNX Runtime

Reuses the same library as layout detection — skip this step if it's
already set up.

- **macOS (Homebrew, CPU only)**: `brew install onnxruntime`
- **macOS with CoreML**: download from
  <https://github.com/microsoft/onnxruntime/releases> (`onnxruntime-osx-arm64-*.tgz`
  or `onnxruntime-osx-x64-*.tgz`)
- **Ubuntu/Debian**: `apt install libonnxruntime-dev`
- **Custom path**: `export ORT_LIB_PATH=/path/to/libonnxruntime.dylib`

## 2. Obtain PaddleOCR ONNX models

The recommended starting point is PP-OCRv5 mobile (det + rec) for Latin
scripts, which covers FR/EN/DE plus other European languages.

```bash
# Create the cache directory
mkdir -p ~/.pdf2md/models/ocr

# Export from PaddleOCR's PaddlePaddle models using paddle2onnx.
# See: https://github.com/PaddlePaddle/PaddleOCR/blob/main/docs/infer_deploy/paddle2onnx.md
#
# Alternatively, use any community-hosted PP-OCRv5 / PP-OCRv4 ONNX pair.
# The files must be named exactly:
cp /path/to/exported/det.onnx  ~/.pdf2md/models/ocr/det.onnx
cp /path/to/exported/rec.onnx  ~/.pdf2md/models/ocr/rec.onnx
```

## 3. (Optional) Use a custom dictionary

The recognition model output indices map to characters via a dictionary.
`pdf2md` ships with the PaddleOCR Latin dictionary embedded, which matches
PP-OCRv5 *Latin* recognition. If your model uses a different dictionary
(Chinese, Cyrillic, etc.), save it one character per line to:

```bash
~/.pdf2md/models/ocr/dict.txt
```

Or pass `--ocr-dict-path` at runtime.

## 4. Verify

```bash
./pdf2md --debug scanned.pdf out.md
# Look for: "ocr: PaddleRecognizer ready"
```

If you see `ocr: det model not found` or similar, re-check step 2.

## Environment variables

| Variable                   | Purpose                              |
| -------------------------- | ------------------------------------ |
| `PDF2MD_OCR_DET_PATH`      | Override the detection model path    |
| `PDF2MD_OCR_REC_PATH`      | Override the recognition model path  |
| `PDF2MD_OCR_DICT_PATH`     | Override the character dictionary    |
| `ORT_LIB_PATH`             | Override the ONNX Runtime library    |

CLI flags take precedence: `--ocr-det-model-path`, `--ocr-rec-model-path`,
`--ocr-dict-path`, `--runtime-path`.

## Disabling OCR

Pass `--no-ocr` to skip OCR entirely — image-only pages will be left empty
in the output. OCR only runs on pages that the custom parser returns as
text-less, so enabling it has zero cost on normal text-bearing PDFs.
