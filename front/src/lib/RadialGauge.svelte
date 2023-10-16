<script lang="ts">
  import { tweened, spring } from 'svelte/motion';
  import { backInOut } from 'svelte/easing';
  import { arc as d3arc } from 'd3-shape';
  import { scaleLinear } from 'd3-scale';
  import * as d3 from "d3"

  export let debug = false;
  export let name:string =""
  export let height:number = 120;
  export let width:number = 150;

  export let value:number = 0;
  export let roundValue:number = 2
  export let unitSymbol = ""
  export let min:number = 0;
  export let max:number = 100;

  let startAngle:number = -120;
  let endAngle:number = 120;
  let innerRadius:number = 50;
  let outerRadius:number = 60;
  let cornerRadius:number = 10;

  let showTextSvgCenter = true;
  let showTextArcCenter = false;
  let showTextArcBottom = false;
  let showTextArcCentroid = false;


  function myScale(value: number) {
    if(value < min ){
      min = Math.floor(value-1)
      //console.log(`[new-min: ${min}, max: ${max}] / val: ${value}`)
    }
    let valueAngle = (d3.scaleLinear().domain([min, max]).range([startAngle, endAngle]).clamp(true))(value);
    if (debug) {
      console.log(name)
      console.log(`[min: ${min}, max: ${max}] / val: ${value}`)
      console.log(`angle:  ${valueAngle}`)
    }
    setColor();
    return valueAngle
  }

  $: valueAngle = myScale(value)

  $: arc = d3arc()
    .innerRadius(innerRadius)
    .outerRadius(outerRadius)
    .startAngle(startAngle * Math.PI / 180)
    .endAngle(valueAngle * Math.PI / 180)
    .cornerRadius(cornerRadius);

  $: trackArc = d3arc()
    .innerRadius(innerRadius)
    .outerRadius(outerRadius)
    .startAngle(startAngle * Math.PI / 180)
    .endAngle(endAngle * Math.PI / 180)
    .cornerRadius(cornerRadius);

  $: trackArcCentroid = trackArc.centroid()
  // $: console.log(trackArcCentroid)

  let trackArcEl
  $: boundingBox = trackArc && trackArcEl ? trackArcEl.getBBox() : {};
  //$: console.log(boundingBox)

  $: textArcCenterOffset = {
    x: (outerRadius - (boundingBox.width  / 2)),
    // x: 0,
    y: (outerRadius - (boundingBox.height  / 2)) * -1
  }
  // $: console.log(textArcCenterOffset)

  $: textArcBottomOffset = {
    x: (outerRadius - (boundingBox.width  / 2)),
    // x: 0,
    y: (outerRadius - (boundingBox.height)) * -1
  }
  // $: console.log(textArcBottomOffset)

  const colors = [
    "rgb(153, 178, 213)",
    "rgb(41, 78, 125)",
    "rgb(72,127,122)",
    "rgb(108,146,114)",
    "rgb(191, 168, 114)",
    "rgb(183, 145, 99)",
    "rgb(179, 121, 87)",
    "rgb(166, 77, 78)",
    "rgb(158, 37, 75)",
    "rgb(70, 17, 38)"]

  const setColor = () => {
    //let i = Math.round((colors.length - min) * value / max)
    let i = Math.floor(d3.scaleLinear().domain([min,max]).range([0, 10])(value))
    //let i = Math.floor((colors.length * value)/ (max-min))
    if (debug) {
      console.log(`color idx: ${i}`)
    }
    color = colors[i]
  }
  let color
  //setColor();


</script>

<svg width={width} height={height} xmlns="http://www.w3.org/2000/svg">
  <path d={trackArc()} transform="translate({width/2}, {height/1.75})" class="track" bind:this={trackArcEl} />
  <path d={arc()} transform="translate({width/2}, {height/1.75})" style="fill: {color}" />
  <text class="range" x={width*.15} y={height*.92}>{min}</text>
  <text class="range" x={width*.8} y={height*.92}>{max}</text>
  {#if showTextSvgCenter}
    <text transform="translate({width/2}, {height/1.75})" dy={16}>
      {value.toFixed(roundValue)}{unitSymbol}
    </text>
  {/if}

  {#if showTextArcCenter}
    <text x={textArcCenterOffset.x} y={textArcCenterOffset.y} transform="translate({width/2}, {height/2})" dy={16}>
      {value}{unitSymbol}
    </text>
  {/if}

  {#if showTextArcBottom}
    <text x={textArcBottomOffset.x} y={textArcBottomOffset.y} transform="translate({width/2}, {height/2})" dy={0}>
      {value}{unitSymbol}
    </text>
  {/if}

  {#if showTextArcCentroid}
    <text x={trackArcCentroid[0]} y={trackArcCentroid[1]} transform="translate({width/2}, {height/2})" dy={16}>
      {value}{unitSymbol}
    </text>
  {/if}

  <defs>
    <linearGradient id="fillGradient" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="rgb(229, 238, 254)"/>
      <stop offset="10%"   stop-color="rgb(153, 178, 213)"/>
      <stop offset="20%" stop-color="rgb(41, 78, 125)"/>
      <stop offset="30%"   stop-color="rgb(117, 147, 135)"/>
      <stop offset="40%" stop-color="rgb(191, 168, 114)"/>
      <stop offset="50%"   stop-color="rgb(183, 145, 99)"/>
      <stop offset="60%" stop-color="rgb(179, 121, 87)"/>
      <stop offset="70%"   stop-color="rgb(166, 77, 78)"/>
      <stop offset="80%" stop-color="rgb(158, 37, 75)"/>
      <stop offset="90%"   stop-color="rgb(70, 17, 38)"/>
    </linearGradient>
  </defs>
</svg>
{#if debug}
<div>
  <label for="">Value:</label>
  <input type="range" {min} {max} bind:value={value} on:input={setColor}/>
  <input type="number" value={value} on:change={setColor}/>
</div>
  {/if}

<style>
    svg {
        /*border: 1px solid #000;*/
    }
    .track {
        stroke: hsla(0, 0%, 0%, .8);
        stroke-width: 1px;
        fill: none;
    }

    text {
        fill: #000;
        font-size: 2rem;
        text-anchor: middle;
    }

    .range {
        font-size: 1rem;
    }

    input[type=number] {
        width: 72px;
    }
</style>