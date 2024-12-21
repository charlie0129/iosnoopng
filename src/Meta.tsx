import { useEffect, useState } from 'react'
import req from './request.ts'
import utils from './utils.ts';
import { Link } from 'react-router';

interface statType {
  exec: string;
  writes: number;
  reads: number;
  lastUpdated: Date;
}

function MetaPage() {
  const [stats, setStats] = useState<statType[]>([])
  const [events, setEvents] = useState<number>(0)
  const [sortBy, setSortBy] = useState<string>('writes')
  const [autoRefresh, setAutoRefresh] = useState<boolean>(true)

  function sortStats(stats: statType[], sortBy: string): statType[] {
    const statsCopy = [...stats]
    return statsCopy.sort((a, b) => {
      if (sortBy === 'writes') {
        return b.writes - a.writes
      } else if (sortBy === 'reads') {
        return b.reads - a.reads
      } else {
        return a.exec.localeCompare(b.exec)
      }
    })
  }

  async function refreshStats() {
    const stats = await req.getStats()
    setEvents(stats.events)
    let statsArray: statType[] = []
    for (const [exec, stat] of Object.entries(stats.processStats)) {
      statsArray.push({ exec, writes: stat.writes, reads: stat.reads, lastUpdated: new Date(stat.lastUpdated) })
    }
    setStats(sortStats(statsArray, sortBy))
  }

  utils.useInterval(() => {
    if (autoRefresh) {
      refreshStats()
    }
  }, autoRefresh ? 2000 : null);

  useEffect(() => {
    refreshStats()
  }, [])

  useEffect(() => {
    setStats(sortStats(stats, sortBy))
  }, [sortBy])

  return (
    <>
      <p>
        <span>Events: {events}</span>&nbsp;&nbsp;&nbsp;
        <span>Writes: {utils.humanFileSize(stats.reduce((acc, stat) => acc + stat.writes, 0))}</span>&nbsp;&nbsp;&nbsp;
        <span>Reads: {utils.humanFileSize(stats.reduce((acc, stat) => acc + stat.reads, 0))}</span>&nbsp;&nbsp;&nbsp;
        <span onClick={() => { setAutoRefresh(!autoRefresh) }} style={{ cursor: "pointer", color: autoRefresh ? "green" : "black" }}>Auto Refresh {autoRefresh ? "On" : "Off"}</span>
      </p>
      <table id="stats-table">
        <thead>
          <tr>
            <th onClick={() => { setSortBy("exec") }}>Process Name&nbsp;&nbsp;</th>
            <th onClick={() => { setSortBy("writes") }}>Data Written&nbsp;&nbsp;</th>
            <th onClick={() => { setSortBy("reads") }}>Data Read&nbsp;&nbsp;</th>
            <th>Last Updated</th>
          </tr>
        </thead>

        <tbody>
          {
            stats.slice(0, 1000).map((stat: statType) => (
              <tr key={stat.exec}>
                <td><Link to={`/${stat.exec}`}>{stat.exec}</Link></td>
                <td>{utils.humanFileSize(stat.writes)}</td>
                <td>{utils.humanFileSize(stat.reads)}</td>
                <td>{utils.timeSince(stat.lastUpdated)}</td>
              </tr>
            ))
          }
        </tbody>
      </table>
      {
        stats.length >= 1000 ?
          <p>Displaying 1000 of {stats.length} entries.</p>
          : null
      }
    </>
  )
}

export default MetaPage
