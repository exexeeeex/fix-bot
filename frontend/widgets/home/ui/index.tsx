import { RegisterMyAbsenseButton } from "@/features/register-my-absence"
import { Button } from "@/shared/components/ui/button"

export const Home: React.FC = () => {
    return (
        <section className="bg-transparent">
            <RegisterMyAbsenseButton userId={"1123"}/>
            <Button variant={'default'}>Я студент</Button>
            <Button >Я преподаватель</Button>
        </section>
    )
}