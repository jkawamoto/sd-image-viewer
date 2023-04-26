import {KeyboardEvent, useCallback} from "react";
import {TextInput} from "@mantine/core";
import {useInputState} from "@mantine/hooks";

type Props = {
  onSearch: (query: string) => void
}

function Query({onSearch}: Props) {
  const [query, setQuery] = useInputState("")
  const onKeyUp = useCallback((event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") {
      event.preventDefault()
      onSearch(query)
    }
  }, [query, onSearch])

  return (
    <TextInput placeholder="Search keywords" value={query} onChange={setQuery} onKeyUp={onKeyUp}/>
  )
}

export default Query
