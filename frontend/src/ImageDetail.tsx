import {Box, Grid, Image, ScrollArea, Stack, Text, Title} from "@mantine/core";
import {Image as ImageInfo} from "./Api.js";
import  dayjs from "dayjs";
import utc from "dayjs/plugin/utc"

dayjs.extend(utc)

function additionalParams(k: string) {
  return k != "id" && k != "prompt" && k != "negative-prompt" && k != "creation-time" && k != "checkpoint"
}


type Props = {
  url: string
  image: ImageInfo
}

function ImageDetail({url, image}: Props) {
  const params = Object.keys(image).filter(additionalParams).map(k => (
    <Box key={k}><Title order={4}>{k}: </Title><Text style={{wordBreak: "break-all"}}>{image[k]}</Text></Box>
  ))
  return (
    <Grid justify="flex-end" align="center">
      <Grid.Col span="auto">
        <Image src={url} alt={image.prompt} radius="md" fit="contain" withPlaceholder height="90vh"/>
      </Grid.Col>
      <Grid.Col span={3} mr="1em">
        <Stack>
          <Box>
            <Title order={3}>ID:</Title>
            <Text style={{wordBreak: "break-all"}}>{image.id}</Text>
          </Box>
          <Box>
            <Title order={3}>Prompt:</Title>
            <Text style={{wordBreak: "break-all"}}>{image.prompt}</Text>
          </Box>
          <Box>
            <Title order={3}>Negative Prompt:</Title>
            <Text style={{wordBreak: "break-all"}}>{image["negative-prompt"]}</Text>
          </Box>
          <Box>
            <Title order={3}>Checkpoint:</Title>
            <Text style={{wordBreak: "break-all"}}>{image.checkpoint}</Text>
          </Box>
          <Box>
            <Title order={3}>Creation Date:</Title>
            <Text style={{wordBreak: "break-all"}}>{dayjs(image["creation-time"]).local().toString()}</Text>
          </Box>
          <Box>
            <ScrollArea h={250}>
              <Stack>
                {params}
              </Stack>
            </ScrollArea>
          </Box>
        </Stack>
      </Grid.Col>
    </Grid>
  )
}


export default ImageDetail
