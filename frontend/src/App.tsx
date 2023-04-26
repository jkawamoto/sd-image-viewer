import {useEffect, useState} from 'react'
import {Api, Image as ImageInfo, Metadata} from "./Api.js";
import {
  AppShell,
  Box,
  Center,
  Footer,
  Grid,
  Header,
  Image,
  Modal,
  Pagination,
  ScrollArea,
  SegmentedControl,
  Stack,
  Text,
  Title,
  UnstyledButton
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
  const [order, setOrder] = useState("desc")

  const api = new Api()
  const getURL = (id: string) => `${api.baseUrl}/${encodeURIComponent(id)}`

  useEffect(() => {
    const fetchImages = async () => {
      const res = await api.getImages({
        query,
        page: page - 1,
        size: size === "small" || size === "medium" || size === "large" ? size : undefined,
        order: order === "asc" ? order : "desc",
      })
      if (res.data.items) {
        setImages(res.data.items)
      }
      setMetadata(res.data.metadata || null)
    }
    fetchImages().catch(console.error)
  }, [page, query, size, order])

  const header = (
    <Header height={{base: 50, md: 70}} p="md" fixed>
      <Grid align="baseline">
        <Grid.Col span="content">
          <SegmentedControl
            value={order}
            onChange={setOrder}
            data={[
              {label: 'Newest', value: 'desc'},
              {label: 'Oldest', value: 'asc'},
            ]}
          />
        </Grid.Col>
        <Grid.Col span="auto">
          <Query onSearch={setQuery}/>
        </Grid.Col>
        <Grid.Col span="content">
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
        </Grid.Col>
      </Grid>
    </Header>
  )
  const footer = (
    <Footer height={{base: 30, md: 50}} fixed>
      <Center mt="0.5em">
        <Pagination total={metadata?.totalPages || 0} onChange={setPage}/>
      </Center>
    </Footer>
  )
  const imageList = (
    images.map((image) => (
      <Grid.Col span={2}>
        <UnstyledButton key={image.id} onClick={() => setImage(image)}>
          <Image src={getURL(image.id)} alt={image.prompt} radius="md" fit="contain" withPlaceholder/>
        </UnstyledButton>
      </Grid.Col>
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
      <Grid justify="flex-start">
        {imageList}
      </Grid>
    </AppShell>
  )
}

export default App
