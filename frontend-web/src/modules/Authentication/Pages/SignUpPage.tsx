import { SignupForm } from "../Components/SignUp/SignUpForm.component"
import { useGetCollegesQuery, useSignupMutation } from "../apis/auth.api"
import toast from "react-hot-toast"

import type { SignupRequest } from "../Interfaces/auth.types"
import type { ApiErrorResponse } from "../../../core/types/api"

export const SignupPage = () => {
	const { data } = useGetCollegesQuery()
	const [signup, { isLoading }] = useSignupMutation()

	const handleSignup = async (payload: SignupRequest) => {
		try {
			await signup(payload).unwrap()
			toast.success("Signup successful. Please login.")
			window.location.href = "/login"
		} catch (err) {
			const error = err as { data?: ApiErrorResponse }
			toast.error(error.data?.error || "Signup failed")
		}
	}

	return (
		<div className="relative h-screen w-screen overflow-hidden">
			{/* Gradient background */}
			<div className="absolute inset-0 bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-900" />
			{/* Glow */}
			<div className="absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.25),transparent_45%)]" />

			{/* Centered content */}
			<div className="relative z-10 flex h-full w-full items-center justify-center px-4">
				<SignupForm
					colleges={data?.data}
					onSubmit={handleSignup}
					loading={isLoading}
				/>
			</div>
		</div>
	)
}


