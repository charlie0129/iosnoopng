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

async function resetStats(): Promise<void> {
    const response = await fetch('/api/stats', { method: 'DELETE' });
    if (!response.ok) {
        throw new Error('Failed to reset stats');
    }
}

async function deleteByProcess(exec: string) {
    const response = await fetch(`/api/stats/${exec}`, { method: 'DELETE' });
    if (!response.ok) {
        throw new Error('Failed to delete process stats');
    }
}

export default { getStats, getProcessStats, resetStats, deleteByProcess };
