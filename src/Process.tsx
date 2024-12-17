import { useEffect, useState } from 'react'
import req from './request.ts'
import utils from './utils.ts';
import { useParams } from 'react-router';

interface statType {
  path: string;
  writes: number;
  reads: number;
}

function ProcessPage() {
  const [stats, setStats] = useState<statType[]>([])
  const [sortBy, setSortBy] = useState<string>('writes')


  let { exec } = useParams();
  exec = exec || "";

  function sortStats(stats: statType[], sortBy: string): statType[] {
    const statsCopy = [...stats]
    return statsCopy.sort((a, b) => {
      if (sortBy === 'writes') {
        return b.writes - a.writes
      } else if (sortBy === 'reads') {
        return b.reads - a.reads
      } else {
        return a.path.localeCompare(b.path)
      }
    })
  }

  useEffect(() => {
    (async () => {
      const stats = await req.getProcessStats(exec)
      let statsArray: statType[] = []
      for (const [path, stat] of Object.entries(stats)) {
        statsArray.push({ path, writes: stat.writes, reads: stat.reads })
      }
      setStats(sortStats(statsArray, sortBy))
    })()
  }, [])

  useEffect(() => {
    setStats(sortStats(stats, sortBy))
  }, [sortBy])

  return (
    <>
      <p>
        <span>Process: {exec}</span>&nbsp;&nbsp;&nbsp;
        <span>Writes: {utils.humanFileSize(stats.reduce((acc, stat) => acc + stat.writes, 0))}</span>&nbsp;&nbsp;&nbsp;
        <span>Reads: {utils.humanFileSize(stats.reduce((acc, stat) => acc + stat.reads, 0))}</span>
      </p>
      <table id="process-table">
        <thead>
          <tr>
            <th onClick={() => { setSortBy("exec") }}>Path</th>
            <th onClick={() => { setSortBy("writes") }}>Data Written</th>
            <th onClick={() => { setSortBy("reads") }}>Data Read</th>
          </tr>
        </thead>

        <tbody>
          {
            stats.slice(0, 1000).map((stat: statType) => (
              <tr key={stat.path}>
                <td>{stat.path}</td>
                <td>{utils.humanFileSize(stat.writes)}</td>
                <td>{utils.humanFileSize(stat.reads)}</td>
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

export default ProcessPage
