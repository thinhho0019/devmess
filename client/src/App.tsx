import { Routes, Route } from "react-router-dom";
import About from "./pages/About";
import NotFound from "./pages/NotFound";
import "./App.css"
import HomeChat from "./pages/HomeChat";
import HomePage from "./pages/Home";
import Login from "./pages/Login";
import AuthSuccess from "./pages/AuthSuccess";
import Register from "./pages/Register";


function App() {

  return (
    <>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/t" element={<HomeChat />} />
        <Route path="/l" element={<Login />} />
        <Route path="/r" element={<Register />} />
        <Route path="/about" element={<About />} />
        <Route path="/auth/success" element={<AuthSuccess />} />
        <Route path="/*" element={<NotFound />} />
      </Routes>
    </>
  )
}

export default App
