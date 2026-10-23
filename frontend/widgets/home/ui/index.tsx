import { RegisterMyAbsenseButton } from "@/features/register-my-absence"

export const Home: React.FC = () => {
    return (
        <section className="bg-transparent">
            <RegisterMyAbsenseButton userId={"1123"}/>
        </section>
    )
}