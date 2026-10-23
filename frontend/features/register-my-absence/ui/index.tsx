import { Icon } from "@/shared/components"
import { Button } from "@maxhub/max-ui"

export const RegisterMyAbsenseButton = (props: {userId: string}) => {
    return <Button  iconBefore={<Icon icon={"plus"}/>}>
        <span >Указать причину отсутствия</span>
    </Button>
}