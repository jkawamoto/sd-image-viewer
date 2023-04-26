import {useEffect, useState} from 'react'
import {Api, Image as ImageInfo, Metadata} from "./Api.js";
import {
  AppShell,
  Center,
  Footer,
  Grid,
  Header,
  Image,
  Modal,
  Pagination,
  SegmentedControl,
  UnstyledButton
} from "@mantine/core";
import Query from "./Query.tsx";
import ImageDetail from "./ImageDetail.tsx";
import {Carousel, Embla, useAnimationOffsetEffect} from "@mantine/carousel";
import {DatePickerInput} from '@mantine/dates';
import dayjs from "dayjs";

const TRANSITION_DURATION = 200;

function App() {
  const [images, setImages] = useState<ImageInfo[]>([])
  const [selectedImage, setSelectedImage] = useState<number | null>(null)
  const [metadata, setMetadata] = useState<Metadata | null>(null)
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState("")
  const [size, setSize] = useState("")
  const [order, setOrder] = useState("desc")
  const [date, setDate] = useState<Date | null>(null);

  const api = new Api()
  const getURL = (id: string) => `${api.baseUrl}/${encodeURIComponent(id)}`

  useEffect(() => {
    const fetchImages = async () => {
      const d = dayjs(date).startOf("day")
      const res = await api.getImages({
        query,
        page: page - 1,
        size: size === "small" || size === "medium" || size === "large" ? size : undefined,
        order: order === "asc" ? order : "desc",
        after: d.toJSON() || undefined,
        before: d.day(d.day() + 1).toJSON() || undefined,
        limit: 12/2*5
      })
      if (res.data.items) {
        setImages(res.data.items)
      }
      setMetadata(res.data.metadata || null)
    }
    fetchImages().catch(console.error)
  }, [page, query, size, order, date])

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
          <DatePickerInput
            onChange={setDate}
            clearable
            placeholder="Filter by date"
            miw="12em"
          />
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
    images.map((image, index) => (
      <Grid.Col span={2} key={image.id}>
        <UnstyledButton onClick={() => setSelectedImage(index)}>
          <Image src={getURL(image.id)} alt={image.prompt} radius="md" fit="scale-down" withPlaceholder/>
        </UnstyledButton>
      </Grid.Col>
    ))
  )
  const detailedImages = images.map((image) => (
    <Carousel.Slide key={image.id}>
      <ImageDetail url={getURL(image.id)} image={image}/>
    </Carousel.Slide>
  ))

  const [embla, setEmbla] = useState<Embla | null>(null);
  useAnimationOffsetEffect(embla, TRANSITION_DURATION);

  return (
    <AppShell padding="md" header={header} footer={footer}>
      <Modal opened={selectedImage !== null} onClose={() => setSelectedImage(null)} fullScreen
             transitionProps={{duration: TRANSITION_DURATION}}>
        <Carousel initialSlide={selectedImage || undefined} draggable={false} getEmblaApi={setEmbla}>
          {detailedImages}
        </Carousel>
      </Modal>
      <Grid justify="flex-start">
        {imageList}
      </Grid>
    </AppShell>
  )
}

export default App
