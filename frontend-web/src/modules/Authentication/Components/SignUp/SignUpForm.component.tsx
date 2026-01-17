import { Button, Input } from "antd"
import toast from "react-hot-toast"
import { useState } from "react"

import { CollegeSelect } from "./CollegeSelect.component"
import { EmailInput } from "./EmailInput.component"
import { PasswordInput } from "./PasswordInput.component"

import { validateCollegeEmail } from "../../utils/validateCollegeEmail"
import { getPasswordStrength } from "../../utils/getPasswordStrength"

import type { College } from "../../../../core/types/college"
import type { SignupRequest } from "../../Interfaces/auth.types"

interface SignupFormProps {
	colleges: College[] | undefined
	loading: boolean
	onSubmit: (payload: SignupRequest) => void
}

export const SignupForm = ({
	colleges,
	onSubmit,
	loading,
}: SignupFormProps) => {
	const [collegeId, setCollegeId] = useState<number | null>(null)
	const [email, setEmail] = useState("")
	const [password, setPassword] = useState("")
	const [name, setName] = useState("")

	const selectedCollege = colleges?.find(
		(c) => c.ID === collegeId
	)

	const passwordStrength = getPasswordStrength(password)

	const handleSubmit = () => {
		if (!collegeId || !selectedCollege) {
			toast.error("Please select college")
			return
		}

		if (!validateCollegeEmail(email, selectedCollege.Domain)) {
			toast.error("Invalid college email")
			return
		}

		if (passwordStrength === "weak") {
			toast.error("Password too weak")
			return
		}

		onSubmit({
			name,
			email,
			password,
			college_id: collegeId,
		})
	}

	return (
		<div
			className="
			w-full max-w-md
			rounded-2xl
			bg-white/60
			backdrop-blur-xl
			border border-white/30
			shadow-2xl
			p-6 sm:p-8
			space-y-5
		"
		>
			{/* Title */}
			<div className="text-center space-y-1">
				<h1 className="text-2xl font-semibold text-gray-900">
					Sign Up
				</h1>
				<p className="text-sm text-gray-700">
					Create your college account
				</p>
			</div>

			<div className="space-y-2">
				<CollegeSelect
					colleges={colleges}
					value={collegeId}
					onChange={setCollegeId}
				/>
				<div>
				<Input
					size="large"
					placeholder="Full name"
					value={name}
					onChange={(e) => setName(e.target.value)}
				/>
				</div>

				<EmailInput value={email} onChange={setEmail} />

				<PasswordInput value={password} onChange={setPassword} />
			</div>

			{/* Submit */}
			<Button
				type="primary"
				size="large"
				block
				loading={loading}
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
				Create account
			</Button>

		</div>
	)

}

