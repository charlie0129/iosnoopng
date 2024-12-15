import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import MetaPage from './Meta.tsx'
import ProcessPage from './Process.tsx'
import {
  BrowserRouter as Router,
  Routes,
  Route,
} from "react-router-dom";

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Router>
      <Routes>
        <Route path="/" element={<MetaPage />}></Route>
        <Route path="/:exec" element={<ProcessPage />}></Route>
      </Routes>
    </Router>
  </StrictMode>,
)
