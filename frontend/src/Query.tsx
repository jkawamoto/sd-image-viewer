import {useState} from "react";
import {TextInput} from "@mantine/core";

type Props = {
  onSearch: (query: string) => void
}

function Query({onSearch}: Props) {
  const [query, setQuery] = useState("")

  return (
    <TextInput placeholder="Search keywords"
               value={query}
               onChange={(event) => setQuery(event.currentTarget.value)}
               onKeyUp={(event) => {
                 if (event.key === "Enter") {
                   event.preventDefault()
                   onSearch(query)
                 }
               }}
    />
  )
}

export default Query
