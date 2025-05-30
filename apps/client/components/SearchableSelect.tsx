import React, { Ref, useState } from "react";
import {
  Combobox,
  ComboboxButton,
  ComboboxInput,
  ComboboxOption,
  ComboboxOptions,
} from "@headlessui/react";

export type Option = { id: string; name: string };

export function SearchableSelect<T extends Option>({
  options,
  onSelect,
  query,
  onQueryChange,
  inputRef,
  additionalContent,
}: {
  options: T[];
  onSelect: (option: T) => void;
  query: string;
  onQueryChange: (query: string) => void;
  inputRef?: Ref<HTMLInputElement>;
  additionalContent?: React.ReactNode;
}) {
  const filteredOptions =
    query === ""
      ? options
      : options.filter((option) => {
          return option.name.toLowerCase().includes(query.toLowerCase());
        });

  return (
    <Combobox
      onChange={(option: T) => {
        if (!option) return;
        onSelect(option);
      }}
      immediate
    >
      <div className="relative w-full">
        <div className="relative w-full cursor-default overflow-hidden rounded-lg bg-white text-left shadow-md focus:outline-none focus-visible:ring-2 focus-visible:ring-white/75 focus-visible:ring-offset-2 focus-visible:ring-offset-teal-300 sm:text-sm">
          <ComboboxInput
            className="w-full border-none py-2 pl-3 pr-10 text-sm leading-5 text-gray-900 focus:ring-0"
            onChange={(event) => {
              onQueryChange(event.target.value);
            }}
            value={query}
            ref={inputRef}
            autoFocus
          />
          <ComboboxButton className="absolute inset-y-0 right-0 flex items-center pr-2">
            {additionalContent}
            <span className="h-5 w-5 text-gray-400" aria-hidden="true">
              ↕️
            </span>
          </ComboboxButton>
        </div>
        <ComboboxOptions
          transition
          className="absolute mt-1 max-h-60 w-full origin-top overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black/5 transition duration-200 ease-out empty:invisible focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0 sm:text-sm"
        >
          {filteredOptions.length === 0 ? (
            <div className="relative cursor-default select-none px-4 py-2 text-gray-700">
              <span className="mr-3">No results found</span>
            </div>
          ) : (
            filteredOptions.map((option) => (
              <ComboboxOption
                key={option.id}
                value={option}
                className="ui-active:bg-teal-600 ui-active:text-white ui-not-active:text-gray-900 relative cursor-default select-none py-2 pl-10 pr-4"
              >
                {option.name}
              </ComboboxOption>
            ))
          )}
        </ComboboxOptions>
      </div>
    </Combobox>
  );
}
