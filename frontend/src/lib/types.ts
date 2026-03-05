export type Transcription = {
    language: string;
    duration: number;
    segments: Segment[];
    id: string;
    translations: Translation[];
    vad_filter?: boolean;
    vad_threshold?: number;
    vad_min_speech_duration_ms?: number;
    vad_min_silence_duration_ms?: number;
}

export type Segment = {
    start: number;
    end: number;
    text: string;
    score: number;
    uuid: string;
    words?: WordData[];
};

type WordData = {
    start: number;
    end: number;
    text: string;
    score: number;
};

type Translation = {
    sourceLanguage: string;
    targetLanguage: string;
    text: string;
    segments: Segment[];
};
