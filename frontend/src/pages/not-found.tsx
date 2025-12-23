import { useEffect } from "react"
import { useNavigate } from "react-router-dom";

export default function NotFoundPage() {
    const navigate = useNavigate();

    useEffect(() => {
        navigate('/');
    }, []);

    return (
        <h1>something went wrong, going home</h1>
    )
}
