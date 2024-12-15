interface RWStat {
    writes: number;
    reads: number;
}

interface Meta {
    [key: string]: RWStat;
}

interface Stats {
    events: number;
    processStatsMeta: Meta;
}

interface ProcessStat {
    [key: string]: Meta;
}

export type { ProcessStat, Stats, Meta, RWStat };