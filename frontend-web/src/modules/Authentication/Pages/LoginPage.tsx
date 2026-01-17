import toast from "react-hot-toast"
import { LoginForm } from "../Components/Login/LoginForm.component"
import { useLoginMutation } from "../apis/auth.api"
import type { FetchBaseQueryError } from "@reduxjs/toolkit/query"
import type { ApiErrorResponse } from "../../../core/types/api"

export const LoginPage = () => {
    const [login, { isLoading }] = useLoginMutation()

    const handleLogin = async (payload: {
        email: string
        password: string
    }) => {
        try {
            await login(payload).unwrap()
            toast.success("Login successful")
            window.location.href = "/"
        } catch (err) {
            const apiErr = err as FetchBaseQueryError & {
                data?: ApiErrorResponse
            }

            toast.error(apiErr.data?.error || "Invalid credentials")
        }
    }

    return (
        <div className="relative h-screen w-screen overflow-hidden">
            {/* Professional background */}
            <div className="absolute inset-0 bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-900" />

            {/* Subtle grain */}
            <div className="absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.08),transparent_50%)]" />

            <div className="relative z-10 flex h-full w-full items-center justify-center px-4">
                <LoginForm onSubmit={handleLogin} loading={isLoading} />
            </div>
        </div>
    )
}
