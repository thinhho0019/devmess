import { Routes, Route  } from "react-router-dom";
import About from "./pages/About";
import NotFound from "./pages/NotFound";
import "./App.css"
import HomeChat from "./pages/HomeChat";
import HomePage from "./pages/Home";  
import Login from "./pages/Login";
 

function App() {
   
  return (
    <>
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/t" element={<HomeChat />} />
      <Route path="/l" element={<Login />} />
      <Route path="/about" element={<About />} />
      <Route path="/*" element={<NotFound />} />
    </Routes>
    </>
  )
}

export default App
