export interface IToast {
    message: string
    type: "error";
}

interface IToastProps {
    toast: IToast;
}

export default function Toast({
    toast,
}: IToastProps) {
    const bg = toast.type == 'error' ? 'bg-red-800' : '';
    return (
        <div className={`ms-auto p-3 max-w-9/10 max-h-28 rounded-md ${bg} wrap-break-word align-middle animate-fade-in-out-right opacity-0 z-10000 pointer-events-none`}>
            {toast.message}
        </div>
    )
}
