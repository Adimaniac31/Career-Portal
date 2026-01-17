import { Input, Progress } from "antd"
import { EyeInvisibleOutlined, EyeOutlined, LockOutlined } from "@ant-design/icons"
import { getPasswordStrength } from "../../utils/getPasswordStrength"

interface PasswordInputProps {
	value: string
	onChange: (value: string) => void
}

export const PasswordInput = ({ value, onChange }: PasswordInputProps) => {
	const strength = getPasswordStrength(value)

	const percent =
		strength === "weak" ? 30 :
		strength === "good" ? 65 :
		strength === "strong" ? 100 : 0

	const status =
		strength === "weak" ? "exception" :
		strength === "good" ? "active" :
		"success"

	return (
		<div>
			<Input.Password
				size="large"
				placeholder="Create password"
				prefix={<LockOutlined />}
				iconRender={visible =>
					visible ? <EyeOutlined /> : <EyeInvisibleOutlined />
				}
				value={value}
				onChange={e => onChange(e.target.value)}
			/>

			{value && (
				<Progress
					percent={percent}
					status={status}
					showInfo={false}
				/>
			)}
		</div>
	)
}
