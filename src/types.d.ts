interface RWStat {
    writes: number;
    reads: number;
    lastUpdated: string;
}

interface Meta {
    [key: string]: RWStat;
}

interface Stats {
    events: number;
    processStats: Meta;
}

interface ProcessStat {
    [key: string]: Meta;
}

export type { ProcessStat, Stats, Meta, RWStat };
