from backends.fasterwhisper import FasterWhisperBackend
from backends.backend import Transcription
from faster_whisper import decode_audio
from models import DeviceType
from typing import Optional
import numpy as np
import io
import os

def convert_audio(file) -> np.ndarray:
        return decode_audio(file, split_stereo=False, sampling_rate=16000)

async def transcribe_from_filename(
    filename: str,
    model_size: int,
    language: Optional[str] = None,
    device: DeviceType = DeviceType.cpu,
    beam_size: int = 5,
    initial_prompt: str = None,
    hotwords: list[str] = None,
    vad_filter: bool = False,
    vad_threshold: Optional[float] = None,
    vad_min_speech_duration_ms: Optional[int] = None,
    vad_min_silence_duration_ms: Optional[int] = None,
) -> Transcription:
    filepath = os.path.join(os.environ["UPLOAD_DIR"], filename)
    if not os.path.exists(filepath):
        raise RuntimeError(f"file not found in {filepath}")
    audio = convert_audio(filepath)
    return await transcribe_audio(audio, model_size, language, device, beam_size, initial_prompt, hotwords, vad_filter, vad_threshold, vad_min_speech_duration_ms, vad_min_silence_duration_ms)

async def transcribe_file(
    file: io.BytesIO,
    model_size: int,
    language: Optional[str] = None,
    device: DeviceType = DeviceType.cpu,
    beam_size: int = 5,
    initial_prompt: str = None,
    hotwords: list[str] = None,
    vad_filter: bool = False,
    vad_threshold: Optional[float] = None,
    vad_min_speech_duration_ms: Optional[int] = None,
    vad_min_silence_duration_ms: Optional[int] = None,
) -> Transcription:
    contents = await file.read()  # async read
    if len(contents) < 150 * 1024 * 1024:  # file is smaller than 150MB
        audio = convert_audio(io.BytesIO(contents))
    else:
        # Save the uploaded file temporarily on disk
        with open(file.filename, 'wb') as f:
            f.write(contents)
        # Check if file exists
        if not os.path.exists(file.filename):
            raise RuntimeError(f"file not found in {file.filename}")
        # Corrected to use the function in this file
        audio = convert_audio(file.filename)
        os.remove(file.filename)
    return await transcribe_audio(audio, model_size, language, device, beam_size, initial_prompt, hotwords, vad_filter, vad_threshold, vad_min_speech_duration_ms, vad_min_silence_duration_ms)

async def transcribe_audio(
    audio: np.ndarray,
    model_size: int,
    language: Optional[str] = None,
    device: DeviceType = DeviceType.cpu,
    beam_size: int = 5,
    initial_prompt: str = None,
    hotwords: list[str] = None,
    vad_filter: bool = False,
    vad_threshold: Optional[float] = None,
    vad_min_speech_duration_ms: Optional[int] = None,
    vad_min_silence_duration_ms: Optional[int] = None,
) -> Transcription:
    if language == "auto":
        language = None
    # Load the model
    model = FasterWhisperBackend(model_size=model_size, device=device)
    model.get_model()
    model.load()
    # Transcribe the file
    return model.transcribe(
        audio,
        silent=True,
        language=language,
        beam_size=beam_size,
        initial_prompt=initial_prompt,
        hotwords=hotwords,
        vad_filter=vad_filter,
        vad_threshold=vad_threshold,
        vad_min_speech_duration_ms=vad_min_speech_duration_ms,
        vad_min_silence_duration_ms=vad_min_silence_duration_ms,
    )
