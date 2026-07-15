import Toast, { IToast } from "./toast";

interface IToastManagerProps {
    toast?: IToast;
}

export default function ToastManager({
    toast,
}: IToastManagerProps) {
    
    return (
        <div className="absolute right-3 bottom-3">
        {
            toast &&
            <Toast toast={toast} />
        }
        </div>
    )
}
