import {useEffect, useState} from 'react'
import {Api, Image as ImageInfo, Metadata} from "./Api.js";
import {
  AppShell,
  Box,
  Center,
  Flex,
  Footer,
  Grid,
  Header,
  Image,
  Modal,
  Pagination,
  ScrollArea,
  Stack,
  Text,
  Title,
  UnstyledButton, SegmentedControl
} from "@mantine/core";
import Query from "./Query.tsx";

function additionalParams(k: string) {
  return k != "id" && k != "prompt" && k != "negative-prompt" && k != "creation-time" && k != "checkpoint"
}


function App() {
  const [images, setImages] = useState<ImageInfo[]>([])
  const [image, setImage] = useState<ImageInfo | null>(null)
  const [metadata, setMetadata] = useState<Metadata | null>(null)
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState("")
  const [size, setSize] = useState("")

  const api = new Api()
  const getURL = (id: string) => `${api.baseUrl}/${encodeURIComponent(id)}`

  useEffect(() => {
    const fetchImages = async () => {
      const res = await api.getImages({
        query,
        page: page - 1,
        size: size === "small" || size === "medium" || size === "large" ? size : undefined,
      })
      if (res.data.items) {
        setImages(res.data.items)
      }
      setMetadata(res.data.metadata || null)
    }
    fetchImages().catch(console.error)
  }, [page, query, size])

  const header = (
    <Header height={{base: 50, md: 70}} p="md" fixed>
      <Flex justify="space-between">
        <Text>SD Image Viewer</Text>
        <Query onSearch={setQuery}/>
        <SegmentedControl
          value={size}
          onChange={setSize}
          data={[
            {label: 'All', value: ''},
            {label: 'Small', value: 'small'},
            {label: 'Medium', value: 'medium'},
            {label: 'Large', value: 'large'},
          ]}
        />
      </Flex>
    </Header>
  )
  const footer = (
    <Footer height={{base: 30, md: 50}} fixed>
      <Center>
        <Pagination total={metadata?.totalPages || 0} onChange={setPage}/>
      </Center>
    </Footer>
  )
  const imageList = (
    images.map((image) => (
      <UnstyledButton key={image.id} onClick={() => setImage(image)}>
        <Image src={getURL(image.id)} alt={image.prompt} radius="md" width={512 / 2} height={768 / 2}
               fit="contain" withPlaceholder/>
      </UnstyledButton>
    ))
  )
  let modal
  if (image) {
    const params = Object.keys(image).filter(additionalParams).map(k => (
      <Box key={k}><Title order={4}>{k}: </Title><Text>{image[k]}</Text></Box>
    ))
    modal = (
      <Modal opened={true} onClose={() => setImage(null)} fullScreen>
        <Grid justify="center" align="center">
          <Grid.Col span={8}>
            <Image src={getURL(image.id)} alt={image.prompt} radius="md" fit="contain" withPlaceholder width="65vw"
                   height="90vh"/>
          </Grid.Col>
          <Grid.Col span={4}>
            <Stack>
              <Box>
                <Title order={3}>Prompt:</Title>
                <Text>{image.prompt}</Text>
              </Box>
              <Box>
                <Title order={3}>Negative Prompt:</Title>
                <Text>{image["negative-prompt"]}</Text>
              </Box>
              <Box>
                <Title order={3}>Checkpoint:</Title>
                <Text>{image.checkpoint}</Text>
              </Box>
              <Box>
                <Title order={3}>Creation Date:</Title>
                <Text>{image["creation-time"]}</Text>
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
      </Modal>
    )
  }

  return (
    <AppShell padding="md" header={header} footer={footer}>
      {modal}
      <Flex wrap="wrap" justify="center" gap={{base: 'sm', sm: 'lg'}}>
        {imageList}
      </Flex>
    </AppShell>
  )
}

export default App
