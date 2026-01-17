import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { Toaster } from "react-hot-toast"
import "antd/dist/reset.css"

import { SignupPage } from "./modules/Authentication/Pages/SignUpPage"
import { LoginPage } from "./modules/Authentication/Pages/LoginPage"

function App() {
	return (
		<BrowserRouter>
			<Toaster position="top-right" />

			<Routes>
				<Route path="/signup" element={<SignupPage />} />
				<Route path="/login" element={<LoginPage />} />

				{/* default */}
				<Route path="*" element={<Navigate to="/login" />} />
			</Routes>
		</BrowserRouter>
	)
}

export default App


