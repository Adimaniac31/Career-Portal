import { Select } from "antd"
import type { College } from "../../../../core/types/college"

interface CollegeSelectProps {
	colleges: College[] | undefined
	value: number | null
	onChange: (value: number) => void
}

export const CollegeSelect = ({
	colleges,
	value,
	onChange,
}: CollegeSelectProps) => {
	return (
		<div>
		<Select<number>
			placeholder="Select your college"
			className="w-full"
			size="large"
			value={value}
			onChange={onChange}
			options={colleges?.map(c => ({
				value: c.ID,
				label: c.Name,
			}))}
		/>
		</div>
	)
}