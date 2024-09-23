import React, { useState } from "react";
import { Dialog, DialogPanel, DialogTitle } from "@headlessui/react";

export type Props = {
  trigger?: (onClick: () => void) => React.ReactNode;
  title: string;
  body: (close: () => void) => React.ReactNode;
  isOpen?: boolean;
  setIsOpen?: (isOpen: boolean) => void;
};

export function Modal({ isOpen, setIsOpen, trigger, title, body }: Props) {
  const [internalIsOpen, internalSetIsOpen] = useState(false);

  isOpen = isOpen ?? internalIsOpen;
  setIsOpen = setIsOpen ?? internalSetIsOpen;

  return (
    <>
      {trigger && trigger(() => setIsOpen(true))}
      <Dialog
        open={isOpen}
        onClose={() => setIsOpen(false)}
        className="relative z-50"
      >
        <div className="fixed inset-0 bg-black/30" aria-hidden="true" />
        <div className="fixed inset-0 flex w-screen items-center justify-center p-4">
          <DialogPanel className="mx-auto max-w-sm rounded bg-white px-4 py-3">
            <DialogTitle className="mb-2 font-bold">{title}</DialogTitle>
            {body(() => setIsOpen(false))}
          </DialogPanel>
        </div>
      </Dialog>
    </>
  );
}
