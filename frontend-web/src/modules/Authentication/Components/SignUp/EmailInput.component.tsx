import { Input } from "antd"
import { MailOutlined } from "@ant-design/icons"

interface EmailInputProps {
	value: string
	onChange: (value: string) => void
}

export const EmailInput = ({ value, onChange } : EmailInputProps) => {
	return (
		<div>
		<Input
			size="large"
			placeholder="Enter college email"
			prefix={<MailOutlined />}
			value={value}
			onChange={e => onChange(e.target.value)}
		/>
		</div>
	)
}
