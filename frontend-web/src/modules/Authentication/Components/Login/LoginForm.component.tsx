import { Button, Input } from "antd"
import { MailOutlined, LockOutlined } from "@ant-design/icons"
import { useState } from "react"

interface LoginFormProps {
	loading: boolean
	onSubmit: (payload: {
		email: string
		password: string
	}) => void
}

export const LoginForm = ({ loading, onSubmit }: LoginFormProps) => {
	const [email, setEmail] = useState("")
	const [password, setPassword] = useState("")
	const BACKEND_URL = import.meta.env.VITE_BACKEND_URL

	const handleSubmit = () => {
		onSubmit({ email, password })
	}

	return (
		<div
			className="
				w-full max-w-md
				rounded-2xl
				bg-white/70
				backdrop-blur-2xl
				border border-white/30
				ring-1 ring-white/20
				shadow-2xl
				p-6 sm:p-8
				space-y-6
			"
		>
			{/* Title */}
			<div className="text-center space-y-1">
				<h1 className="text-2xl font-semibold text-gray-900">
					Welcome back
				</h1>
				<p className="text-sm text-gray-700">
					Login to your college portal
				</p>
			</div>

			{/* Fields */}
			<div className="space-y-4">
				<Input
					size="large"
					placeholder="College email"
					prefix={<MailOutlined />}
					value={email}
					onChange={(e) => setEmail(e.target.value)}
				/>

				<Input.Password
					size="large"
					placeholder="Password"
					prefix={<LockOutlined />}
					value={password}
					onChange={(e) => setPassword(e.target.value)}
				/>
			</div>

			{/* Submit */}
			<Button
				type="primary"
				size="large"
				block
				loading={loading}
				disabled={!email || !password}
				onClick={handleSubmit}
				className="
					!h-12
					!rounded-xl
					!bg-indigo-600
					hover:!bg-indigo-700
					active:scale-[0.97]
					transition-all
				"
			>
				Log in
			</Button>

			{/* Divider */}
			<div className="flex items-center gap-3">
				<div className="h-px flex-1 bg-gray-300/60" />
				<span className="text-xs text-gray-500">OR</span>
				<div className="h-px flex-1 bg-gray-300/60" />
			</div>

			{/* SSO */}
			<Button
				size="large"
				block
				onClick={() => {
					window.location.href = `${BACKEND_URL}/api/auth/sso/login`
				}}
			>
				Continue with College SSO
			</Button>
		</div>
	)
}
