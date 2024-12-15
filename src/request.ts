import type { Meta, Stats } from './types.d.ts';

async function getStats(): Promise<Stats> {
    const response = await fetch('/api/stats');
    if (!response.ok) {
        throw new Error('Failed to fetch stats');
    }
    return response.json();
}

async function getProcessStats(exec: string): Promise<Meta> {
    const response = await fetch(`/api/stats/${exec}`);
    if (!response.ok) {
        throw new Error('Failed to fetch process stats');
    }
    return response.json();
}

export default { getStats, getProcessStats };
